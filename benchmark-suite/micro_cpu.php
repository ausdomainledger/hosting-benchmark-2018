<?php

function execute() {
  $c = 0;
  $t = 0;
  $l = 0;
  $n = 0;

  for ($c = 3; $c < 10000; $c++) {
    $t = floor(sqrt($c));
    for ($l = 2; $l <= $t; $l++) {
      if ($c % $l === 0) {
        break;
      }
    }
    if ($l > $t) {
      $n++;
    }
  }

  return $n;
}

ob_start();

// Warming CPU caches and JIT
for ($i = 0; $i < 250; $i++) {
  execute();
}

// Measured
$start = microtime(true);
$executions = 0;
// Reading PHP's microtime.c, it's clear that microtime uses the gettimeofday syscall, but actual perf recordings
// with php-*-dbgsym indicate that the call is essentially negligible (0.00% of sampled overhead) and is very likely going via
// vDSO anyway. So we are more than OK with calling it in a tight loop.
while (microtime(true) - $start < 30) {
  execute();
  $executions++;
}

$dt = microtime(true) - $start;

echo json_encode([
  "micro_cpu.php" => [
    "executions" => $executions,
    "duration" => $dt
  ]
], JSON_PRETTY_PRINT), PHP_EOL;

ob_end_flush();

?>