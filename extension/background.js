chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    console.log("background")
    console.log(request.message)
    if (request.message === "open_max_url") {
        fullURL = "http://" + request.url;
        chrome.tabs.create({
            "url": fullURL,
            "active": false
        });
    }
});


// window.onload = function () {
//   console.log("im alive ws://localhost:1337/ws");
//   var socket = new WebSocket("ws://localhost:1337/ws");
//   socket.onopen = function () {
//     console.log("Socket is open");
//     socket.send(JSON.stringify({ fuckyeah: "fuckyeah" }));
//   };
//   socket.onmessage = function (e) {
//     console.log("Got some shit:" + e.data);
//     chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
//       chrome.tabs.executeScript(
//         tabs[0].id,
//         { code: 'alert("wow")' },
//         (resulsts) => console.log(resulsts)
//       );
//     });
//   };
//   socket.onclose = function () {
//     console.log("Socket closed");
//   };
// };
//

// chrome.tabs.onActivated.addListener(function(info) {
//     console.log(info);
// });
//
// chrome.tabs.onUpdated.addListener(function(info) {
//     console.log(info);
// });
//
//
// chrome.extension.onConnect.addListener(function(port){
//     port.onMessage.addListener(factory);
// });
//
// function factory(obj){
//     console.log(obj)
// }
//
//
