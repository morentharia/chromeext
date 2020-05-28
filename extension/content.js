function DOMEval( code, doc ) {
    doc = doc || document;

    var script = doc.createElement( "script" );

    script.text = code;
    doc.body.appendChild( script ).parentNode.removeChild( script );
}

chrome.runtime.onMessage.addListener(function (request, _, _) {
  if (request.message === "eval") {
    console.log(request.highlighted_code);
    // TODO: make eval and eval_dom
    // let result = new Function(request.code)();
    let result = DOMEval(request.code)
    chrome.runtime.sendMessage({
      message: "eval_done",
      result: result || null,
    });
  }
});
