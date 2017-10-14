<?php
$dir = @$_REQUEST['dir'];
if (!$dir) {
    $dir = '.';
}
if (file_exists($dir . '/index.json')) {
    $json = json_decode(file_get_contents($dir . '/index.json'), true, 512, JSON_UNESCAPED_UNICODE);
} else {
    die();
}
?>
<!DOCTYPE HTML>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
    <link rel="stylesheet" href="./node_modules/photoswipe/dist/photoswipe.css">
    <link rel="stylesheet" href="./node_modules/photoswipe/dist/default-skin/default-skin.css">
    <link rel="stylesheet" href="./node_modules/minireset.css/minireset.min.css">
    <link rel="stylesheet" href="./flex.css">
    <style>
        h2 {
            font-size: 2em;
        }
        h2 a {
            display: inline-block;
            border: 1px solid black;
            border-radius: 10px;
            padding: 5px 10px;
            text-decoration: none;
            color: black;
            background-color: #ccc;
        }
        h2.breadcrumbs a:last-child {
            background-color: #fff;
            font-weight: bold;
        }
        h2.folders {
            display: flex;
            flex-wrap: wrap;
        }
        h2.folders a {
            font-size: 0.7em;
            padding: 10px 10px;
            display: block;
            width: 100%;
            max-width: 400px;
            background-color: #ccf;
        }
        body { font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; font-size: 1.02em; }
    </style>
</head>
<body>
<h2 class="breadcrumbs"><a href="index.php">Fotos</a><?php
    if ($json['file'] != '') {
        $breadcrumbs = explode('/', $json['file']);
        $breadcrumb = '';
        foreach ($breadcrumbs as $folder) {
            $breadcrumb .= $folder;
            echo "<a href=\"index.php?dir=$breadcrumb\">$folder</a>";
            $breadcrumb .= '/';
        }
    }
    ?> <?=$json['date']?></h2>
<?php if (count($json['folders'])) {?>
<h2 class="folders">
    <?php
        foreach ($json['folders'] as $folder) {
            if ($json['file'] == '') {
                echo "<a href=\"index.php?dir=$folder\">$folder</a>";
            } else {
                echo "<a href=\"index.php?dir=${json['file']}/$folder\">$folder</a>";
            }
        }
    ?>
</h2>
<?php } ?>
<?php if (count($json['images'])) {?>
<div id="images" class="flex-images">
<?php
    foreach ($json['images'] as $image) {
        echo "<a class=\"item\"
    href=\"./${json['file']}/${image['file']}.small.jpg\"
    data-sub-html=\"${image['name']} - ${image['date']} - ${image['width']} x ${image['height']}\"
    data-w=\"${image['width']}\"
    data-h=\"${image['height']}\">
    <img src=\"./${json['file']}/${image['file']}.thumb.jpg\"><div class=\"over\">${image['name']}</div>
</a>\n";
    }
?>
</div>
<script src="./node_modules/javascript-flex-images/flex-images.min.js"></script>
<script>
    new flexImages({selector: '#images', rowHeight: 200});
</script>
<script src="./node_modules/photoswipe/dist/photoswipe.min.js"></script>
<script src="./node_modules/photoswipe/dist/photoswipe-ui-default.min.js"></script>
<script>
    var items = [
<?php
    foreach ($json['images'] as $image) {
        echo "{
    src: './${json['file']}/${image['file']}.small.jpg',
    w: ${image['width']},
    h: ${image['height']},
    title:'${image['name']} - ${image['date']} - ${image['width']} x ${image['height']}',
    msrc:'./${json['file']}/${image['file']}.thumb.jpg'
},";
    }
?>
    ];
    document.getElementById('images').onclick = function(event) {
        var elm = event.target;
        while (elm.className !== 'item') {
            elm = elm.parentElement;
        }
        var pswpElement = document.querySelectorAll('.pswp')[0];
        var options = {
            index: Array.prototype.indexOf.call(elm.parentElement.children, elm),
            getThumbBoundsFn: function(index) {
                var pageYScroll = window.pageYOffset || document.documentElement.scrollTop;
                // get position of element relative to viewport
                var rect = elm.parentElement.children[index].getBoundingClientRect();
                return {x:rect.left, y:rect.top + pageYScroll, w:rect.width};
            },
            showAnimationDuration: 100,
            hideAnimationDuration: 200
        };

        var gallery = new PhotoSwipe(pswpElement, PhotoSwipeUI_Default, items, options);
        gallery.init();
        return false;
    };
</script>

<!-- Root element of PhotoSwipe. Must have class pswp. -->
<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">

    <!-- Background of PhotoSwipe.
         It's a separate element as animating opacity is faster than rgba(). -->
    <div class="pswp__bg"></div>

    <!-- Slides wrapper with overflow:hidden. -->
    <div class="pswp__scroll-wrap">

        <!-- Container that holds slides.
            PhotoSwipe keeps only 3 of them in the DOM to save memory.
            Don't modify these 3 pswp__item elements, data is added later on. -->
        <div class="pswp__container">
            <div class="pswp__item"></div>
            <div class="pswp__item"></div>
            <div class="pswp__item"></div>
        </div>

        <!-- Default (PhotoSwipeUI_Default) interface on top of sliding area. Can be changed. -->
        <div class="pswp__ui pswp__ui--hidden">

            <div class="pswp__top-bar">

                <!--  Controls are self-explanatory. Order can be changed. -->

                <div class="pswp__counter"></div>

                <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>

                <button class="pswp__button pswp__button--share" title="Share"></button>

                <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>

                <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>

                <!-- Preloader demo http://codepen.io/dimsemenov/pen/yyBWoR -->
                <!-- element will get class pswp__preloader--active when preloader is running -->
                <div class="pswp__preloader">
                    <div class="pswp__preloader__icn">
                        <div class="pswp__preloader__cut">
                            <div class="pswp__preloader__donut"></div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
                <div class="pswp__share-tooltip"></div>
            </div>

            <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)">
            </button>

            <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)">
            </button>

            <div class="pswp__caption">
                <div class="pswp__caption__center"></div>
            </div>

        </div>

    </div>

</div>
<?php } ?>
</body>
</html>