function DOMEval(code, doc) {
  doc = doc || document;
  let script = doc.createElement("script");
  script.text = code;
  doc.body.appendChild(script).parentNode.removeChild(script);
}

chrome.runtime.onMessage.addListener(function (request, _, _) {
  if (request.message_type === "dom_eval") {
    console.log(request.highlighted_code);
    let result = DOMEval(request.code);
    chrome.runtime.sendMessage({
      _id: request._id,
      message_type: "dom_eval_done",
      result: result || null,
    });
  } else if (request.message_type === "eval") {
    console.log(request.highlighted_code);
    let result = new Function(request.code)();
    chrome.runtime.sendMessage({
      _id: request._id,
      message_type: "eval_done",
      result: result || null,
    });
  } else if (request.message_type === "ping") {
    chrome.runtime.sendMessage({
      _id: request._id,
      message_type: "pong",
    });
  } else {
    console.warn("unknown message" + request.message);
  }
});
