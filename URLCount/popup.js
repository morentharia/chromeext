// $(function () {
  // Send a message to content.js to fetch all the top domains
  chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
    var activeTab = tabs[0];
    chrome.tabs.sendMessage(activeTab.id, { message: "fetch_top_domains" });
  });

  // Add a listener to handle the response from content.js
  chrome.runtime.onMessage.addListener(function (request, _, _) {
    console.log(request);
    if (request.message === "all_urls_fetched") {
      console.log("all_urls_fetched");
      var urlWithMaxLinks;
      var maxLinks = 0;

      for (var key in request.data) {
        if (request.data.hasOwnProperty(key)) {
          $("#tabs_table tr:last").after(
            "<tr><td>" +
              key +
              "</td>" +
              "<td>" +
              request.data[key] +
              "</td></tr>"
          );

          if (request.data[key] > maxLinks) {
            maxLinks = request.data[key];
            urlWithMaxLinks = key;
          }
        }
      }

      if (maxLinks != 0) {
        chrome.runtime.sendMessage({
          message: "open_max_url",
          url: urlWithMaxLinks,
        });
      }
    }

    if (request.message === "eval") {
      chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, request);
      });
    }

    if (request.message === "something_completed") {
      console.log("popup something_completed");
      chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, request);
      });
    }
  });
// });
