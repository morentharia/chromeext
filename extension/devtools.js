chrome.devtools.panels.create("Razdva36", null, "/devpanel.html", function(panel) {});

chrome.devtools.inspectedWindow.eval("HAHA()",
                                     { useContentScriptContext: true });

