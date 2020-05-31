var urls = [];
for(var i = document.links.length; i-- > 0;) {
    // if(document.links[i].hostname !== location.hostname)
    urls.push(document.links[i].href);
}

for (var j = 0, len = urls.length; j < len; j++) {
    console.log(urls[j]);
    break
}
return urls
