const shell = require('electron').shell;

$(function(){
  $('.open-in-browser').click((event) => {
    event.preventDefault();
    shell.openExternal(event.target.href);
  });
});
