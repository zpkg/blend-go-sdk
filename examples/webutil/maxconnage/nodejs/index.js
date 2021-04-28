/**
 * Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
 * Use of this source code is governed by a MIT license that can be found in the LICENSE file.
 */
const http = require('http');

const httpAgent = new http.Agent({ 
  maxSockets: 32,
  keepAlive: true,
});

async function doTest() {
  for (;;) {
    await doRequest()
  }
}

const newRequest = (url, options) => {
  let lifecycle = {
    didStart: false,
    didData: false,
    didError: false,
    didEnd: false,
  }
  return new Promise((resolve, reject) => {
    let req = http.request(url, options, (res) => {
      lifecycle.didStart = true;
      var body = [];
      res.on('data', (chunk) => {
        lifecycle.didData = true;
        body.push(chunk);
      });
      res.on('end', () => {
        lifecycle.didEnd = true;
        resolve({
          res,
          body,
          lifecycle
        });
      });
    });
    req.on('error', (err) => {
      lifecycle.didError = true;
      reject({req, err, lifecycle});
    });
    req.end();
  });
}
 
let requests = 0;
let successes = 0;
let errors = 0;

async function doRequest() {
  requests++;
  return newRequest("http://127.0.0.1:8081", {
    agent: httpAgent,
  }).then(() => {
    successes++
  }).catch((res) => { 
    if (res.didStart && res.didData && res.didError) {
      errors++;
    }
  })
}

async function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

setTimeout(_ => {
  console.log("quitting");
  console.log("requests:", requests);
  console.log("successes:", successes);
  console.log("errors:", errors);
  process.exit(0)
}, 60 * 1000)
doTest();