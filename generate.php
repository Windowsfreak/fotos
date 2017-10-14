<?php
define('MAGICK_EXECUTABLE', '"C:\Program Files\ImageMagick-7.0.7-Q16\magick"');
define('SQL_DSN', 'mysql:host=localhost;dbname=database;charset=utf8');
define('SQL_USER', 'admin');
define('SQL_PASS', 'password');
ob_end_flush();
/*
   CREATE TABLE `database`.`files` (
       `path` VARCHAR(512) NOT NULL,
       `parent` VARCHAR(512) NOT NULL,
       `type` TINYINT UNSIGNED NOT NULL,
       `date` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       PRIMARY KEY (`path`),
       INDEX `parents` (`parent`),
       INDEX `date` (`date`)
   ) ENGINE = InnoDB;
 */
function createPath($tardir) {
    @mkdir($tardir, 0700, true);
}
function toDate($date) {
    preg_match_all('/\d+/', $date, $tmp, PREG_SET_ORDER, 0);
    if (count($tmp) < 2) {
        return $date;
    }
    if (count($tmp) > 2 && count($tmp) < 6) {
        if (!isset($tmp[3])) {
            $tmp[3] = ['00'];
        }
        if (!isset($tmp[4])) {
            $tmp[4] = ['00'];
        }
        if (!isset($tmp[5])) {
            $tmp[5] = ['00'];
        }
    }
    return $tmp[2][0].'.'.$tmp[1][0].'.'.$tmp[0][0].' '.$tmp[3][0].':'.$tmp[4][0].':'.$tmp[5][0];
}
function identifyImage($srcdir, $infile) {
    $srcpath = str_replace('\\', '/', realpath($srcdir) . '\\' . $infile);
    $path = MAGICK_EXECUTABLE . " identify -format \"%t#%w#%h#%[EXIF:DateTimeOriginal]\" \"${srcpath}\"";
    $val = explode('#', `$path`);
    return [
        'file' => $infile,
        'name' => $val[0],
        'width' => $val[1],
        'height' => $val[2],
        'date' => $val[3]
    ];
}
function resizeImage($srcdir, $infile, $tardir, $outfile, $width, $height) {
    createPath($tardir);
    $srcpath = str_replace('\\', '/', realpath($srcdir) . '\\' . $infile);
    $tarpath = str_replace('\\', '/', realpath($tardir) . '\\' . $outfile);
    $path = MAGICK_EXECUTABLE . " convert \"${srcpath}\" -thumbnail ${width}x${height}^> -filter lanczos2sharp -quality 60 \"${tarpath}\" 2>&1";
    `$path`;
}
$db = new PDO(SQL_DSN, SQL_USER, SQL_PASS, array(PDO::ATTR_EMULATE_PREPARES => false, PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION));
function obtainData($basedir, $minutes) {
    global $db;
    $stmt = $db->prepare('SELECT * FROM files WHERE (path = ? OR parent = ?) AND date > DATE_SUB(NOW(), INTERVAL ? MINUTE)');
    $stmt->execute(array($basedir, $basedir, $minutes));
    $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);
    $files = [];
    foreach ($rows as $row) {
        $files[$row['path']] = [
            'type' => $row['type'],
            'date' => $row['date']
        ];
    }
    return $files;
}
function storeData($basedir, $type) {
    global $db;
    $stmt = $db->prepare('INSERT INTO files (path, parent, type, date) VALUES(?, ?, ?, NOW()) ON DUPLICATE KEY UPDATE type = ?, date=NOW()');
    $parent = dirname($basedir);
    $stmt->execute(array($basedir, $parent, $type, $type));
}
function checkFolder($srcdir, $tardir, $sqldir, $name) {
    $contents = obtainData($sqldir, 525600);
    echo $sqldir . '<br>';
    flush();
    if (isset($contents[$sqldir])) {
        return;
    }
    $scan = scandir($srcdir);
    $folders = [];
    $files = [];
    $images = [];
    $dates = [];
    foreach ($scan as $file) {
        if ($file[0] == '.') {
            continue;
        }
        $filename = "$srcdir/$file";
        if (is_dir($filename)) {
            if ($file == 'snapshot') continue;
            $folders[] = $file;
        } else {
            if (strpos($file, '.small.jpg') !== false || strpos($file, '.thumb.jpg') !== false) continue;
            if (!in_array(strtolower(pathinfo($filename, PATHINFO_EXTENSION)), ['jpg', 'jpeg'])) continue;
            $files[] = $file;
        }
    }
    natsort($folders);
    natsort($files);
    foreach ($folders as $folder) {
        $filename = "$srcdir/$folder";
        if (!isset($contents["$sqldir/$folder"]) || $contents["$sqldir/$folder"]['type'] != 0) {
            checkFolder($filename, "$tardir/$folder", "$sqldir/$folder", $folder);
        }
    }
    foreach ($files as $file) {
        $filename = "$srcdir/$file";
        $info = identifyImage($srcdir, $file);
        if (!isset($contents["$sqldir/$file"]) || $contents["$sqldir/$file"]['type'] != 1) {
            resizeImage($srcdir, $file, $tardir, $file.'.small.jpg', 800, 800);
            resizeImage($tardir, $file.'.small.jpg', $tardir, $file.'.thumb.jpg', 200, 200);
            storeData("$sqldir/$file", 1);
        }
        if ($info['date'] != '') {
            $dates[] = $info['date'];
            $info['date'] = toDate($info['date']);
        }
        $info['file'] = "$file";
        $images[] = $info;
    }
    $date = '';
    if (count($dates) > 0) {
        natsort($dates);
        $firstDate = toDate($dates[0]);
        $lastDate = toDate(end($dates));
        reset($dates);
        if ($firstDate == $lastDate) {
            $date = $firstDate;
        } else {
            $date = $firstDate . ' - ' . $lastDate;
        }
    }
    $fp = @fopen("$tardir/index.json", 'w');
    if ($fp) {
        fwrite($fp, json_encode([
            'name' => $name,
            'file' => $sqldir,
            'folders' => $folders,
            'images' => $images,
            'date' => $date
        ], JSON_UNESCAPED_UNICODE));
        fclose($fp);
    }
    storeData($sqldir, 0);
}
ini_set('max_execution_time', 30000);
checkFolder('Z:/Fotos', 'Fotos', 'Fotos', 'Fotos');