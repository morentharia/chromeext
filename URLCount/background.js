console.log("im alive ws://localhost:1337/ws");

var socket = new WebSocket("ws://localhost:1337/ws");

socket.onopen = function () {
  console.log("Socket is open");
  // socket.send(JSON.stringify({ fuckyeah: "fuckyeah" }));
};

socket.onmessage = function (e) {
  console.log("Got some shit:");
  console.log(JSON.parse(e.data));
  // chrome.runtime.sendMessage(JSON.parse(e.data));

  console.log("send something_completed");
  chrome.runtime.sendMessage({
    message: "something_completed",
    data: {
      subject: "Loading",
      content: "Just completed!",
    },
  });

  // chrome.tabs.create({"url": "https://fuck", "active": false});
};

socket.onclose = function () {
  console.log("Socket closed");
};


chrome.runtime.onMessage.addListener(
  function(request, _, _) {
    console.log(request.message)
    if( request.message === "open_max_url" ) {
      fullURL = "http://" + request.url;
      chrome.tabs.create({"url": fullURL, "active": false});
    }
  }
);

  // chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
  //   chrome.tabs.executeScript(
  //     tabs[0].id,
  //     { code: 'alert("wow")' },
  //     (resulsts) => console.log(resulsts)
  //   );
  // });
