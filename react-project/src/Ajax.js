import axios from 'axios';

const Ajax = function(request){
  if(!request.params){
    request.params = {};
  }
  if(!request.scope){
    request.scope = this;
  }
  if(!request.failure){
    request.failure = function(){}
  }

  axios
  .post(request.url, request.params) //, {Post:Name}
  .then((resp)=>{
    console.log('resp');
    console.log(resp);
    if(resp.data && resp.data.success && resp.data.success === true){
      request.success.call(request.scope, resp.data)
    }else{
      var err
      if(resp.data.error){
        err = resp.data.error;
      }else{
        err = "Ajax error! Check the console!";
      }
      console.log(err);
      console.log(resp);
      request.failure.call(request.scope)
      alert(err);
    }
  })
  .catch((resp)=>{
    request.failure.call(request.scope, arguments)
  });
}

export default Ajax;