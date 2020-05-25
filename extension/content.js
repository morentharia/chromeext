chrome.runtime.onMessage.addListener(function (request, _, _) {
  if (request.message === "eval") {
    console.log(request.highlighted_code);
    let result = new Function(request.code)();
    chrome.runtime.sendMessage({
      message: "eval_done",
      result: result || null,
    });
  }
});
