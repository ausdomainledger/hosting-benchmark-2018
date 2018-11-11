<?php

// PHP does not expose the right syscalls for us
// to be able to meaningfully measure much.
// At most, we get fflush(3), which can just result
// in dirty/unwritten pages in the kernel.
// Database and WordPress benchmarks likely reflect on
// IO performance in a more meaningful way.

$WORK_DIR = dirname(__FILE__) . "/io_latency_work";
// To be more representative for shared web hosting, 
// we will prefer many little inodes.
$NUM_FILES = 2000;
$FILE_SIZE = 32 * 1024;
$BLOCK_SIZE = 4 * 1024;

function cleanup() {
  global $WORK_DIR;
  foreach (glob("$WORK_DIR/*") as $node) {
    unlink($node);
  }
  if (is_dir($WORK_DIR))
    rmdir($WORK_DIR);
}

cleanup();
if (!mkdir($WORK_DIR)) {
  die("mkdir");
}

ob_start();

$open_duration = 0;
$seq_write_duration = 0;

for ($i = 0; $i < $NUM_FILES; $i++) {
  $noise = random_bytes($BLOCK_SIZE);

  // inode creation
  $start = microtime(true);
  $fh = fopen("$WORK_DIR/data_$i", "w");
  if (!$fh) {
    die("fopen");
  }
  $open_duration += microtime(true) - $start;

  // sequential write
  $start = microtime(true);
  for ($written = 0; $written < $FILE_SIZE; $written += $BLOCK_SIZE) {
    if (!fwrite($fh, $noise)) {
      die("fwrite");
    }
  }
  fflush($fh);
  $seq_write_duration += microtime(true) - $start;
  fclose($fh);
}

// crappy mixed random read and write workload for 15s
$random_rw_count = 0;
$random_rw_duration = 0;
$start = microtime(true);
while (true) {
  if (($random_rw_duration = microtime(true) - $start) >= 15) {
    break;
  }
  // choose random inode
  $r = mt_rand(0, $NUM_FILES-1);
  $fh = fopen("$WORK_DIR/data_$r", "a+");
  if (!$fh) {
    die("fopen");
  }

  // read random 4k block
  $r = mt_rand(0, ($FILE_SIZE / $BLOCK_SIZE) - 1);
  if (fseek($fh, $r * $BLOCK_SIZE) === -1) {
    die("fseek");
  }
  if(!fread($fh, $BLOCK_SIZE)) {
    die("fread");
  }

  // write random 4k block
  $r = mt_rand(0, ($FILE_SIZE / $BLOCK_SIZE) - 1);
  if (fseek($fh, $r * $BLOCK_SIZE) === -1) {
    die("fseek");
  }
  if (!fwrite($fh, $noise)) {
    die("fwrite");
  }

  fflush($fh);
  fclose($fh);

  $random_rw_count++;
}

cleanup();

echo json_encode([
  "micro_io.php" => [
    "open_duration" => $open_duration,
    "seq_write_duration" => $seq_write_duration,
    "random_rw_duration" => $random_rw_duration,
    "random_rw_count" => $random_rw_count
  ],
], JSON_PRETTY_PRINT), PHP_EOL;

ob_end_flush();
?>