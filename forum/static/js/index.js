"use strict";

document.addEventListener("DOMContentLoaded", init);

function init(){    
    let btnGenApi = document.getElementById('gapik');
    btnGenApi.onclick = genAPIKey;
}

function genAPIKey(event){
    var myAlert = document.getElementById('apikalert');
    var opt = {
        animation: true,
        autohide: true,
        delay: 300000
      };
    var bsAlert = new bootstrap.Toast(myAlert, opt)    
    
    var url="/getapikey";
    var options={
        method: 'GET',
        headers: {
        Accept: 'application/json',
        'Content-Type': 'application/x-www-form-urlencoded'
        }        
    };
    let response = fetch(url, options);
    response.then((response)=>{        
        let result = response.json();         
        result.then((data)=>{                        
            var apiKey = data.APIKey; 
            var apikMsg = document.getElementById("apikalerttxt");
            apikMsg.innerText="APIKey:\n"+apiKey;
            bsAlert.show();          
        });
    });    
}