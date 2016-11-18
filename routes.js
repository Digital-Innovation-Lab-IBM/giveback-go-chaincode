var express = require('express')
var router = express.Router()

var ibc
var chaincode

module.exports.setup = function(sdk, cc){
  ibc = sdk;
  chaincode = cc;
}

module.exports.isSetup = function(){
  return ibc && chaincode
}

router.get('/user/:userId', function(req, res){
  chaincode.query.read([req.params.userId], function(e, data){
    if(e){
      console.log("Error " + e)
      res.status(400);
      res.send("Error " + e);
    }
    else if(!data){
      console.log("Error - Data not found for some reason?")
      res.status(400);
      res.send("Error - Data not found for some reason?");
    }
    else{
      console.log(data)
      res.status(200);
      res.send(data);
    }
  })
})

router.post('/trade/:senderId/:receiverId/:amount', function(req, res){
  chaincode.query.read([req.params.senderId], function(e, data){
    if(e || !data){
      console.log("Error - Sender user doesnt exist")
      res.status(400);
      res.send("Error - Sender user doesnt exist");
    }
    else{
      chaincode.query.read([req.params.receiverId], function(e, data){
        if(e || !data){
          console.log("Error - Receiver user doesnt exist")
          res.status(400);
          res.send("Error - Receiver user doesnt exist");
        }
        else{
          chaincode.invoke.set_user([req.params.senderId, req.params.amount, req.params.receiverId], function(e, data){
            if(e){
              console.log("Error " + e)
              res.status(400);
              res.send("Error " + e);
            }
            else if(!data){
              console.log("Error - Data not found for some reason?")
              res.status(400);
              res.send("Error - Data not found for some reason?");
            }
            else{
              console.log(data)
              res.status(200);
              res.send(data);
            }
          })
        }
      })
    }
  })
})

router.post('/createAccount/:username', function(req, res){
  chaincode.query.read([req.params.username], function(e, data){
    if(e || data){
      console.log("User already exists")
      res.status(400);
      res.send("User already exists");
    }
    else{
      chaincode.invoke.createAccount([req.params.username], function(e, data){
        if(e){
          console.log("Error " + e)
          res.status(400);
          res.send("Error " + e);
        }
        else if(!data){
          console.log("Error - Data not found for some reason?")
          res.status(400);
          res.send("Error - Data not found for some reason?");
        }
        else{
          console.log(data)
          res.status(200);
          res.send(data);
        }
      })
    }
  })
})

module.exports.router = router
