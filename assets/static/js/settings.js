var settingsChanges = {
  RemovedFeeds: [],
  AddedFeeds: [],
  UpdatesEvery: "-1",
  ExpiresAfter: "-1",
}

var str1 = '<tr><td>';
var str2 = '</td><td><button type="button" class="btn btn-outline-warning"><i class="fas fa-plus-circle"></i></button></td></tr>';

$(function(){
  $('.btn-settings-add').click(function(){
    value = $(this).parent().parent().children().first().children().first().val();
    if (!validator.isURL(value)) {
      $('div.modal-body').prepend('<div class="alert alert-danger alert-dismissible fade show" role="alert"><strong>Added feed is not a url.<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>');
      return;
    }
    settingsChanges.AddedFeeds.push(value);
    $(this).parent().parent().children().first().children().first().val('');
    $(this).parent().parent().parent().before(str1 + value + str2);
  });
  $('.btn-settings-delete').click(function(){
    value = $(this).parent().parent().children().first().text();
    settingsChanges.RemovedFeeds.push(value);
    $(this).parent().parent().remove();
  });
  $('#btnSettingsSave').click(function(){
    ue = $('#updatesEvery').val()
    if (ue) {
      settingsChanges.UpdatesEvery = ue
    }
    ea = $('#expiresAfter').val()
    if (ea) {
      settingsChanges.ExpiresAfter = ea
    }
    $('div.modal-body').prepend('<div class="alert alert-success alert-dismissible fade show" role="alert"><strong>Saving data!<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>');
    $.post( "/save", JSON.stringify(settingsChanges))
      .done(function(data){
        window.location.reload(true);
      });
  });
});
