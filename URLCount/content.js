// content.js

chrome.runtime.onMessage.addListener(function (request, _, _) {
  if (request.message === "something_completed") {
      console.log("fuck yeah something_completed");
      chrome.runtime.sendMessage({ message: "eval_done"});
  }

  if (request.message === "eval") {
      console.log("eval")
      console.log(request);
      console.log(request);
      console.log(request);
      eval(request.message.code)
      chrome.runtime.sendMessage({ message: "eval_done"});
  }

  if (request.message === "fetch_top_domains") {
    var urlHash = {},
      links = document.links;
    for (var i = 0; i < links.length; i++) {
      var domain = links[i].href.split("/")[2];
      if (urlHash[domain]) {
        urlHash[domain] = urlHash[domain] + 1;
      } else {
        urlHash[domain] = 1;
      }
    }
    chrome.runtime.sendMessage({ message: "all_urls_fetched", data: urlHash });
  }
});
