//content.js
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    if (request.message === "fetch_top_domains") {
        // Handle the message
        chrome.runtime.sendMessage({ "message": "all_urls_fetched" });
    }
});
