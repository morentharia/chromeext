function connect() {
  console.log("im alive");

  var socket = new WebSocket("ws://localhost:1337/ws");

  socket.onopen = function () {
    console.log("Socket is open");
  };

  socket.onmessage = function (e) {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      var activeTab = tabs[0];
      if (typeof activeTab === "undefined") {
        console.log("no active tab");
        return;
      }
      chrome.tabs.sendMessage(activeTab.id, JSON.parse(e.data));
      //
      //
      //
      // console.log(">>>>>>>>>" + JSON.parse(e.data).code);
      // chrome.tabs.executeScript({
      //   code:  JSON.parse(e.data).code,
      // });
      //
      //
      // console.log("hahahahah");
      // console.log("hahahahah");
      // console.log("hahahahah");
      //
      //
    });
  };

  socket.onclose = function () {
    console.log("Socket closed reconnect after 1sec");
    setTimeout(function () { connect(); }, 1000);
  };

  socket.onerror = function (err) {
    console.error("Socket encountered error: ", err.message, "Closing socket");
    socket.close();
  };

  chrome.runtime.onMessage.addListener(function (request, _, _) {
    console.log(request);
    socket.send(JSON.stringify(request));
  });
}

connect();
