(function(w) {
  w.galleries = {
    items: {},
    item_count: 0,
    first_load: true,
    selected: 0,
    base_url: "",

    buildViewURL: function(id) {
      return galleries.base_url + galleries.items[id].id
    },

    refreshItems: function() {
      $.getJSON(feed, function(data) {
        galleries.item_count = data.count;
        galleries.items = data.objects;
        galleries.redrawThumbs();
      });
    },

    resizeCanvas: function() {
      var node = $('#gallery_image > a > img')[0];

      // This is so we can get the base dimensions of the image.
      var image = new Image();
      image.src = galleries.buildViewURL(galleries.selected);

      var max_height = w.innerHeight - 96;
      var max_width = Math.max(w.innerWidth - 150, 550);

      var ar = image.height / image.width;
      var scale = Math.min(1, (max_width/image.width), (max_height/image.height));

      var new_width = image.width * scale;
      var new_height = image.width * ar * scale;

      node.height = new_height;
      node.width  = new_width;
      $(node).attr("style", "position:absolute;top:74px;left:10px;"); // Remove the temporary max-width and max-height
    },

    redrawSelected: function() {
      var parent = $('#gallery_image > a');
      parent.empty();

      var max_height = w.innerHeight - 96;
      var max_width = Math.max(w.innerWidth - 150, 550);

      parent.append("<img style='max-width:"+max_width+"px;max-height:"+max_height+"px;'/>");
      $('#gallery_image > a > img').imagesLoaded(galleries.resizeCanvas);
      parent.find('img')[0].src = galleries.buildViewURL(galleries.selected);
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