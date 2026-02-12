(function(w) {
  w.galleries = {
    items: {},
    item_count: 0,
    first_load: true,
    selected: 0,
    base_url: "",

    buildViewURL: function(id) {
      return galleries.base_url + galleries.items[id].id;
    },

    isVideo: function(id) {
      var item = galleries.items[id];
      if (!item) return false;

      // Check mime type if available
      if (item.mime_type) {
        return item.mime_type.indexOf('video/') === 0;
      }

      // Fallback to URL extension
      var url = galleries.buildViewURL(id);
      var ext = url.split('.').pop().toLowerCase().split('?')[0];
      var videoExts = ['mp4', 'webm', 'mov', 'ogv'];
      return videoExts.indexOf(ext) !== -1;
    },

    refreshItems: function() {
      $.getJSON(feed, function(data) {
        galleries.item_count = data.count;
        galleries.items = data.objects;
        galleries.redrawThumbs();
      });
    },

    resizeCanvas: function() {
      var node = $('#gallery_image > a > img')[0] || $('#gallery_image > a > video')[0];
      if (!node) return;

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
      }

      $(node).attr("style", "position:absolute;top:74px;left:10px;");
    },

    redrawSelected: function() {
      var parent = $('#gallery_image > a');
      parent.empty();

      var max_height = w.innerHeight - 96;
      var max_width = Math.max(w.innerWidth - 150, 550);

      if (galleries.isVideo(galleries.selected)) {
        var videoUrl = galleries.buildViewURL(galleries.selected);
        parent.append(
          "<video style='max-width:" + max_width + "px;max-height:" + max_height + "px;' " +
          "controls autoplay loop muted playsinline>" +
          "<source src='" + videoUrl + "' type='video/mp4'>" +
          "</video>"
        );
        $('#gallery_image > a > video').on('loadedmetadata', galleries.resizeCanvas);
      } else {
        parent.append("<img style='max-width:" + max_width + "px;max-height:" + max_height + "px;'/>");
        $('#gallery_image > a > img').imagesLoaded(galleries.resizeCanvas);
        parent.find('img')[0].src = galleries.buildViewURL(galleries.selected);
      }
    },

    redrawThumbs: function() {
      $('#gallery_thumbs').fadeOut(function(){
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
      $('#thumb_'+id).removeClass('unselected').addClass('selected');

      if (keyboard) {
        $("#gallery_sidebar").stop().scrollTo("#thumb_"+id, 500, { offset: -5 });
      }

      $('#gallery_image > a').attr('href', galleries.buildViewURL(galleries.selected));
      galleries.redrawSelected();

      return false;
    },

    init: function() {
      galleries.base_url = base_img_url + '/';
      galleries.refreshItems();

      $(w).bind('resize', galleries.resizeCanvas);
      $(document).keydown(function(e){
        switch(e.keyCode){
          case 37: // left
          case 38: // up
            if (galleries.selected > 0) {
              galleries.view((galleries.selected-1), true);
              return false;
            }
            break;

          case 39: // right
          case 40: // down
            if (galleries.selected < (galleries.item_count-1)) {
              galleries.view((galleries.selected+1), true);
              return false;
            }
            break;
        }
      })
    },
  }
  $(w.document).ready(galleries.init);
})(window);