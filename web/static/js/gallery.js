(function(w) {
  w.galleries = {
    items: {},
    item_count: 0,
    first_load: true,
    selected: 0,
    base_url: "",
    is_resizing: false,
    resize_timeout: null,

    buildViewURL: function(id) {
      var item = galleries.items[id];
      if (!item) return "#";
      return galleries.base_url + item.id;
    },

    isVideo: function(id) {
      var item = galleries.items[id];
      if (!item) return false;
      return item.mime_type.indexOf('video/') === 0;
    },

    isAudio: function(id) {
      var item = galleries.items[id];
      if (!item) return false;
      return item.mime_type.indexOf('audio/') === 0;
    },

    refreshItems: function() {
      $.getJSON(feed, function(data) {
        galleries.item_count = data.count;
        galleries.items = data.objects;
        galleries.redrawThumbs();
      });
    },

    resizeCanvas: function() {
      if (galleries.is_resizing) return;
      galleries.is_resizing = true;

      var node = $('#gallery_image > a > img')[0] ||
                 $('#gallery_image > a > video')[0] ||
                 $('#gallery_image > a > audio')[0];

      if (!node) {
        galleries.is_resizing = false;
        return;
      }

      if (node.tagName === 'AUDIO') {
        var max_width = Math.max(w.innerWidth - 150, 550);
        $(node).css({
          'position': 'absolute',
          'top': '74px',
          'left': '10px',
          'width': max_width + 'px'
        });
        galleries.is_resizing = false;
        return;
      }

      var nodeWidth, nodeHeight;

      if (node.tagName === 'VIDEO') {
        nodeWidth = node.videoWidth || node.offsetWidth;
        nodeHeight = node.videoHeight || node.offsetHeight;
      } else {
        nodeWidth = node.naturalWidth || node.width;
        nodeHeight = node.naturalHeight || node.height;
      }

      var max_height = w.innerHeight - 96;
      var max_width = Math.max(w.innerWidth - 150, 550);

      if (nodeWidth > 0) {
        var ar = nodeHeight / nodeWidth;
        var scale = Math.min(1, (max_width / nodeWidth), (max_height / nodeHeight));
        var new_width = nodeWidth * scale;
        var new_height = nodeWidth * ar * scale;
        node.height = new_height;
        node.width = new_width;
        $(node).attr("style", "position:absolute;top:74px;left:10px;");
      }

      galleries.is_resizing = false;
    },

    redrawSelected: function() {
      var parent = $('#gallery_image > a');
      parent.empty();

      var max_height = w.innerHeight - 96;
      var max_width = Math.max(w.innerWidth - 150, 550);

      if (galleries.isVideo(galleries.selected)) {
        var videoUrl = galleries.buildViewURL(galleries.selected);
        var mimeType = galleries.items[galleries.selected].mime_type || 'video';
        parent.append(
          "<video style='max-width:" + max_width + "px;max-height:" + max_height + "px;' " +
          "controls loop muted playsinline>" +
          "<source src='" + videoUrl + "' type='" + mimeType + "'>" +
          "</video>"
        );
        var $video = $('#gallery_image > a > video');
        $video.one('loadedmetadata', galleries.resizeCanvas);
        $video[0].load();
      } else if (galleries.isAudio(galleries.selected)) {
        var audioUrl = galleries.buildViewURL(galleries.selected);
        var mimeType = galleries.items[galleries.selected].mime_type || 'audio/mpeg';
        parent.append(
          "<audio style='width:" + max_width + "px;' controls>" +
          "<source src='" + audioUrl + "' type='" + mimeType + "'>" +
          "</audio>"
        );
        galleries.resizeCanvas();
      } else {
        parent.append("<img style='max-width:" + max_width + "px;max-height:" + max_height + "px;'/>");
        $('#gallery_image > a > img').imagesLoaded(galleries.resizeCanvas);
        parent.find('img')[0].src = galleries.buildViewURL(galleries.selected);
      }
    },

    handleResize: function() {
      clearTimeout(galleries.resize_timeout);
      galleries.resize_timeout = setTimeout(galleries.resizeCanvas, 100);
    },

    redrawThumbs: function() {
      $('#gallery_thumbs').fadeOut(function() {
        $('#gallery_thumbs').html("");
        $('#galleryThumbs').tmpl(galleries.items).appendTo('#gallery_thumbs');
        $('#gallery_thumbs').fadeIn(function() {
          if (galleries.first_load) {
            galleries.view(0);
            galleries.first_load = false;
          }
        });
      });
    },

    view: function(id, keyboard) {
      if (id == galleries.selected && !galleries.first_load)
        return;

      galleries.selected = id;

      $('#gallery_thumbs > .puush_tile').removeClass('selected').addClass('unselected');
      $('#thumb_' + id).removeClass('unselected').addClass('selected');

      if (keyboard) {
        $("#gallery_sidebar").stop().scrollTo("#thumb_" + id, 500, { offset: -5 });
      }

      $('#gallery_image > a').attr('href', galleries.buildViewURL(galleries.selected));
      galleries.redrawSelected();

      return false;
    },

    init: function() {
      galleries.base_url = base_img_url + '/';
      galleries.refreshItems();

      $(w).bind('resize', galleries.handleResize);
      $(document).keydown(function(e) {
        switch (e.keyCode) {
          case 37:
          case 38:
            if (galleries.selected > 0) {
              galleries.view((galleries.selected - 1), true);
              return false;
            }
            break;

          case 39:
          case 40:
            if (galleries.selected < (galleries.item_count - 1)) {
              galleries.view((galleries.selected + 1), true);
              return false;
            }
            break;
        }
      });
    },
  };
  $(w.document).ready(galleries.init);
})(window);