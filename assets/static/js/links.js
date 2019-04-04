const shell = require('electron').shell;

$(function(){
  $('.open-in-browser').click((event) => {
    event.preventDefault();
    shell.openExternal(event.target.href);
  });
  $('.lime-btn-link').click((event) => {
    $(event.target).addClass('disabled');
    $.post("/read", $(event.target).text());
  });
});

document.addEventListener('astilectron-ready', function() {
  astilectron.onMessage(function(message) {
    if (message === "refresh") {
      window.location.reload(true); 
      return;
    }
  });
})
