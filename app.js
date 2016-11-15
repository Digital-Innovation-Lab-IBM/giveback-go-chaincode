'use strict';

var express = require('express');
var http = require('http');
var app = express();
var bodyParser = require('body-parser');
var fs = require('fs');
var cors = require('cors');

app.use(bodyParser.json());
app.use(bodyParser.urlencoded());

app.options('*', cors());
app.use(cors());

app.set('port', 8080);
app.listen(app.get('port'), '0.0.0.0', function() {
  // print a message when the server starts listening
  console.log("server starting on " + app.get('port'));
});

var Ibc1 = require('ibm-blockchain-js');														//rest based SDK for ibm blockchain
var ibc = new Ibc1();

try{
	var manual = JSON.parse(fs.readFileSync('mycreds.json', 'utf8'));
	var peers = manual.credentials.peers;
	console.log('loading hardcoded peers');
	var users = null;																			//users are only found if security is on
	if(manual.credentials.users) users = manual.credentials.users;
	console.log('loading hardcoded users');
}
catch(e){
	console.log('Error - could not find hardcoded peers/users, this is okay if running in bluemix');
}

function prefer_type1_users(user_array){
	var ret = [];
	for(var i in users){
		if(users[i].enrollId.indexOf('type1') >= 0) {	//gather the type1 users
			ret.push(users[i]);
		}
	}

	if(ret.length === 0) ret = user_array;				//if no users found, just use what we have
	return ret;
}

var options = 	{
	network:{
		peers: peers,																	//lets only use the first peer! since we really don't need any more than 1
		users: prefer_type1_users(users),											//dump the whole thing, sdk will parse for a good one
		options: {
					quiet: true 															//detailed debug messages on/off true/false
				}
	},
	chaincode:{
		zip_url: 'https://github.com/ibm-blockchain/marbles/archive/v2.0.zip',
		unzip_dir: 'marbles-2.0/chaincode',													//subdirectroy name of chaincode after unzipped
		git_url: 'http://gopkg.in/ibm-blockchain/marbles.v2/chaincode',		//GO get http url
	}
};

var chaincode = null;
ibc.load(options, function (err, cc){														//parse/load chaincode, response has chaincode functions!
	if(err != null){
		console.log('! looks like an error loading the chaincode or network, app will fail\n', err);
		if(!process.error) process.error = {type: 'load', msg: err.details};				//if it already exist, keep the last error
	}
	else{
		chaincode = cc;
		// PASS ANY CC OBJECTS TO WS

		// ---- To Deploy or Not to Deploy ---- //
		if(!cc.details.deployed_name || cc.details.deployed_name === ''){					//yes, go deploy
			cc.deploy('init', ['99'], {delay_ms: 30000}, function(e){ 						//delay_ms is milliseconds to wait after deploy for conatiner to start, 50sec recommended
				check_if_deployed(e, 1);
			});
		}
		else{																				//no, already deployed
			console.log('chaincode summary file indicates chaincode has been previously deployed');
			check_if_deployed(null, 1);
		}
	}
});

//loop here, check if chaincode is up and running or not
function check_if_deployed(e, attempt){
	if(e){
		console.log('! looks like a deploy error, holding off on the starting the socket\n', e);
	}
	else if(attempt >= 15){																	//tried many times, lets give up and pass an err msg
		console.log('[preflight check]', attempt, ': failed too many times, giving up');
		var msg = 'chaincode is taking an unusually long time to start. this sounds like a network error, check peer logs';
		cb_deployed(msg);
	}
	else{
		console.log('[preflight check]', attempt, ': testing if chaincode is ready');
		chaincode.query.read(['_marbleindex'], function(err, resp){
			var cc_deployed = false;
			try{
				if(err == null){															//no errors is good, but can't trust that alone
					if(resp === 'null') cc_deployed = true;									//looks alright, brand new, no marbles yet
					else{
						var json = JSON.parse(resp);
						if(json.constructor === Array) cc_deployed = true;					//looks alright, we have marbles
					}
				}
			}
			catch(e){}																		//anything nasty goes here

			// ---- Are We Ready? ---- //
			if(!cc_deployed){
				console.log('[preflight check]', attempt, ': failed, trying again');
				setTimeout(function(){
					check_if_deployed(null, ++attempt);										//no, try again later
				}, 10000);
			}
			else{
				console.log('[preflight check]', attempt, ': success');
        chaincode.query.read(['a'])
				//cb_deployed(null);															//yes, lets go!
			}
		});
	}
}
