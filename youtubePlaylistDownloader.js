// Use http://www.downvids.net/
// Enter playlist
// Copypaste the entire file into the javascript file and wait for the output at the very end.

// After getting output from javascript console pass it into the python file cuz fuck
// trying to make your browser download files with ajax.


// Converts all the "files" on the webpage by calling the onclick function
// in the webpage
function convertAll(){
  var a = $('a').get();

  var b = Array();
  a.forEach(function(data) {
      if(data.id == "search_more"){
          if (data.outerHTML.indexOf("convert(")>-1){
              var start = data.outerHTML.indexOf("convert(");
              var end = data.outerHTML.indexOf("getads()");
              b.push(data.outerHTML.substring(start,end));
          }
      }
  });
  b.forEach(function(data){eval(data)});


  setTimeout(grabAllLinks,20000);
}

// Will grab all links from converted videos. 
// this assumes that convertAll() has finished executing
// and that the webpage itself has finished converting.
// Give ample time before calling grabAllLinks() after convertAll()
function grabAllLinks(){
  var b = Array();
  var a = $("a").get();
  a.forEach(function(data) {
      if(data.id == "search_more"){
          if (data.outerHTML.indexOf("Mp3 conver")>-1){
              console.log("a");
              var start = data.outerHTML.indexOf("href=\"http://")+6;
              var end = data.outerHTML.indexOf(".mp3")+4;
              b.push(data.outerHTML.substring(start,end));
          }
      }
  });
  var urls = "";
  b.forEach(function(data) {
      urls += data + "|"
  });
  console.log(urls);
  console.log("Done. Copy the above output and paste it into the .py file.")
}

convertAll();