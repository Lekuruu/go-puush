var puush_mode = 'view';
var busy = false;
var dialog_visible = false;
var view = "grid";
var checkboxes;
var old_username = "";
var timeout;

function puush_toggle_mode()
{
  if (puush_mode == 'select')
  {
    // exit selection mode
    puush_mode = 'view';
    $('#select-bar').removeClass('select-mode');
    $('.select').fadeOut();
    puush_clear_selection();
    $('#select-button .label').text("Select");
  }
  else if (puush_mode == 'view')
  {
    // switch to selection mode
    $('#select-bar').addClass('select-mode');
    puush_mode = 'select';
    $('.select').fadeIn();
    $('#select-button .label').text("Cancel");
  }
  return false;
}

function toggleEverything(currentCheckbox)
{
    checkboxes.each(function(i, box) { box.checked = currentCheckbox.checked; });
    puush_update_selection();
}

function puush_update_selection()
{
  var selected = puush_get_selected().length;
  var suffix = "puush'd file" + (selected != 1 ? "s" : "");

  if (selected > 0)
  {
    $('#rm-link').text("Delete ("+selected+")");
    $('#mv-link').text("Move ("+selected+")");
    $('#tg-link').text("Tag ("+selected+")");

    $('#del-dialog-msg').html("Deleting <strong>"+selected+"</strong> "+suffix+"...");
  } else {
    $('#rm-link').text("Delete");
    $('#mv-link').text("Move");
    $('#tg-link').text("Tag");
  }
}

// return an array of ids corrosponding to the selected files
function puush_get_selected()
{
  var marked = Array();

  if (view == "grid")
  {
    $('.puush_tile.selected').each(
      function (index, element)
      {
        marked.push(element.id.substring(6));
      }
    );
  } else {
    checkboxes.each( function(i, box) { if ( box.checked ) { marked.push(box.value) } } );
  }

  return marked;
}

function puush_hist_select(id)
{
  if (!busy)
  {
    if (view == "grid")
    {
      $('#puush_'+id).toggleClass("unselected");
      $('#puush_'+id).toggleClass("selected");
    } else {
      $('#tr-'+id).toggleClass("selected");
      var box = $('#tr-'+id+' .td-ck input')[0];
      box.checked = !box.checked;
    }
    puush_update_selection();
  }
  return false;
}

function puush_click(id)
{
  switch (puush_mode)
  {
    case 'view':
      return true;
    break;

    case 'select':
      puush_hist_select(id);
      return false;
    break;
  }
}

function toggle_busy(id)
{
  if (busy)
  {
    $('#'+id).html("");
    busy = false;
  } else {
    busy = true;
    $('#'+id).html("<img src='/img/spinner.gif' />");
  }
}

function puush_be_busy(action)
{
  // close the blinds
  $('#facebox').hide();
  $('#blinds').fadeIn();//function(){$('#throbber').fadeIn();});
  $('#throbber').fadeIn();
  $('#throbber-msg').text(action);

}

function puush_not_busy()
{
  // open the blinds
  $('#blinds').fadeOut();
  $('#throbber').fadeOut();
}

function puush_clear_selection()
{
  if (view == "grid")
  {
    $('.puush_tile.selected').each(
      function (index, element)
      {
        $(element).removeClass('selected');
        $(element).addClass('unselected');
      }
    );
  } else {
    checkboxes.each( function(i,box) { box.checked = false; $('#tr-'+box.value).removeClass("selected"); } );
  }

  puush_update_selection();
  return false;
}

function puush_select_all()
{
  if (view == "grid")
  {
    $('.puush_tile').each(
      function (index, element)
      {
        $(element).removeClass('unselected');
        $(element).addClass('selected');
      }
    );
  } else {
    checkboxes.each( function(i,box) { box.checked = true; $('#tr-'+box.value).addClass("selected"); } );
  }

  puush_update_selection();
  return false;
}

function puush_confirm_delete()
{
  var marked = puush_get_selected();

  if (marked.length > 0)
  {
    $("#blinds").fadeIn();
    $("#del-dialog").fadeIn();
    dialog_visible = true;
  }
  return false;
}

function puush_do_delete()
{
  var marked = puush_get_selected();

  if (marked.length > 0)
  {
    $("#del-dialog").fadeOut();
    puush_be_busy("Deleting " + marked.length + " file(s)");
    $.post("/ajax/delete_upload", { 'i[]': marked },
      function(data)
      {
        window.location.reload();
      }
    );
  }
  return false;
}

function puush_confirm_move()
{
  var marked = puush_get_selected();

  if (marked.length > 0)
  {
    /* TODO: make not ugly */
    puush_be_busy("Preparing to Move");
    $.ajax({
      url: '/ajax/move_dialog/?i=' + marked.join(','),
      success:
        function(data)
        {
          $("#move-dialog-msg").html(data);
          $("#throbber").fadeOut();
          $("#move-dialog").fadeIn();
          dialog_visible = true;
        }
    });
  }
  return false;
}

function puush_do_move(pool)
{
  var marked = puush_get_selected();

  $("#move-dialog").fadeOut();
  puush_be_busy("Moving " + marked.length + " file(s)");

  $.post("/ajax/move_upload", { 'i[]': marked, 'p': pool },
    function(data)
    {
      window.location.reload();
    }
  );

  return false;
}

function puush_change_password()
{
  var cur_pass = document.change_password_form[1].current_password.value;
  var new_pass = document.change_password_form[1].new_password.value;
  var cnf_pass = document.change_password_form[1].confirm_password.value;

  if (new_pass === cnf_pass)
  {
    $.post("/ajax/change_password", { 'c': cur_pass, 'p': new_pass },
      function(data)
      {
        if (data['error'] == true)
        {
          alert(data['message']);
        } else {
          $(document).trigger('close.facebox');
        }
      }
    );
  } else {
    alert("Uh-oh! Those passwords don't match!");
  }
  return false;
}

function set_default_pool(id)
{
  var result  = $('#pool_result');
  var spinner = $('#pool_spinner');
  $.post("/ajax/default_puush_pool", { 'p': id }, function(data) { result.hide(); spinner.fadeIn(); result.html(data); result.fadeIn(); spinner.fadeOut(); result.fadeOut()} );
}

function close_popup()
{
  if (dialog_visible)
  {
    $("#blinds").fadeOut();
    $("#move-dialog").fadeOut();
    $("#del-dialog").fadeOut();
    $("#welcome-dialog").fadeOut();
    $('#tweetbot').fadeOut();
    $('#tweetbot_wrapper').fadeOut();
  }
  dialog_visible = false;
  return false;
}

function close_facebox()
{
  jQuery(document).trigger('close.facebox');
}

function on_keydown(event)
{
  if (event.keyCode == '27')
    close_popup();
}

function search_focus()
{
  var box = document.getElementById('search-box');
  if (box.value == "Search...")
    box.value = "";
}

function search_blur()
{
  var box = document.getElementById('search-box');
  if (box.value == "")
    box.value = "Search...";
}

function new_username_wizard()
{
  // close the blinds
  dialog_visible = true;
  $('#blinds').fadeIn(function(){$('#welcome-dialog').fadeIn()});

}

function stop_asking_about_username()
{
  $.post("/ajax/stopnagging", {'n':'username'}, function(data) {
    if (!data.error)
    {
      close_popup();
    }
  }, "json");
  return false;
}

function puush_check_username()
{
  var username = $('#new_username').val();
  $.post("/ajax/confirm_username", { 'u': username }, function(data) {
    if (data.error)
    {
      $("#username-ok")[0].src = "/img/nok.png";
      $("#username-error-note").text(data.message);
    } else {
      $("#username-ok")[0].src = "/img/ok.png";
      $("#username_confirm_button")[0].disabled = false;
    }
  }, "json");
  return false;
}

function puush_claim_username()
{
  var username = $('#new_username').val();
  $.post("/ajax/claim_username", { 'u': username }, function(data) {
    if (data.error)
    {
      $("#username-ok")[0].src = "/img/nok.png";
      $("#username-ok").fadeIn();
      $("#username-error-note").text(data.message);
    } else {
      window.location.reload();
    }
  }, "json");
  return false;
}

function puush_update_sample_gallery_path()
{
  var new_username = $('#new_username').val();

  if (new_username.length >= 3 && new_username != old_username)
  {
    $("#username-ok").fadeIn();
    $("#username-ok")[0].src = "/img/spinner.gif";
    old_username = new_username;
    $("#sample-gallery-url > span").text($("#new_username").val());
    $("#username-error-note").html("&nbsp;");
    $("#username_confirm_button")[0].disabled = true;
    clearTimeout(timeout);
    timeout = setTimeout(puush_check_username, 500);
  }
}

function help_tb()
{
	$('#facebox').hide();
	$('#blinds').fadeIn();
	$('#tweetbot').nivoSlider({manualAdvance: true, directionNavHide: false, keyboardNav:true, nextOnClick:true, animSpeed:100, effect: 'fade'});
	$('#tweetbot').fadeIn();
	$('#tweetbot_wrapper').fadeIn();
	dialog_visible = true;
}

function request_tb_email()
{
	puush_be_busy('Requesting email...');
	$.post("/ajax/tb_email", function(data) { puush_not_busy(); });
}

function init()
{
  $(document).bind('beforeReveal.facebox', function() { $('#blinds').fadeIn(); });
  $(document).bind('close.facebox', function() { $('#blinds').fadeOut(); });
  $(document).bind('keydown', on_keydown);
  $('#new_username').bind('keydown', puush_update_sample_gallery_path);
  $('#new_username').bind('keyup', puush_update_sample_gallery_path);

  $('a[rel*=facebox]').facebox();
  $("#blinds").bind('click', close_popup);

  if ( $('#puush-select').length == 1 )
    view = "list";

  if (view == "list")
  {
    checkboxes = $('#puush-select input[name=i]');
    checkboxes.bind('click', function(event) { event.target.checked = !event.target.checked; /* fucking lol :( */ puush_update_selection(); });
  }

  $("#search-box").bind('focus', search_focus);
  $("#search-box").bind('blur', search_blur);

  if (new_username_prompt && new_username_nag)
    new_username_wizard();
}

function exportPool(poolId)
{
  if (!confirm("This will export all files in this pool. Depending on the size of your pool, this could take a while. You will receive an email once this operation is complete, with a link to download the archived export."))
    return;

  $.post("/ajax/export_pool", { 'p': poolId }, function(data) {
    switch (parseInt(data))
    {
      case -2:
        alert("There is nothing in this pool to export!");
        break;
      case -1:
        alert("An export is already running. Please wait to receive your email.");
        break;
      default:
        alert("Error has occurred.");
        break;
      case 1:
        alert("Export is in progress. You will receive an email when it has completed.");
        break;
    }
  });
}

$(document).ready(init);