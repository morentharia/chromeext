chrome.tabs.query(
  {
    active: true,
    currentWindow: true,
  },
  function (tabs) {
    var activeTab = tabs[0];
    chrome.tabs.sendMessage(activeTab.id, { message: "fetch_top_domains" });
  }
);

chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    if (request.message === "all_urls_fetched") {
        console.log("all_urls_fetched");
    }
});
