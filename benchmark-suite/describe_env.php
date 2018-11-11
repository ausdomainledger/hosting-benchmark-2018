<?php

// Capture phpinfo
ob_start();
phpinfo();
$phpinfo_capture = ob_get_clean();

// Find the io scheduler for real block devices
$schedulers = [];
foreach (glob("/sys/block/*/queue/scheduler") as $dev) {
  if (strpos($dev, "loop")) {
    continue;
  }
  $schedulers[$dev] = trim(file_get_contents($dev));
}

// Find the MySQL version from the initial wire packet
$mysql_version = '';
$conn = @fsockopen("localhost", 3306, $errno, $errstr, 1);
if ($conn) {
  fread($conn, 4 + 1); // packet header<4> + proto version<1>
  while (($b = fread($conn, 1)) != chr(0)) { // string.NUL
    $mysql_version .= $b;
  }
  fclose($conn);
}

// Extract physical and swap
$meminfo = file_get_contents("/proc/meminfo");
preg_match('/MemTotal:\s+(\d+)/m', $meminfo, $memtotal);
$memtotal = count($memtotal) > 1 ? $memtotal[1] : '';
preg_match('/SwapTotal:\s+(\d+)/m', $meminfo, $swaptotal);
$swaptotal = count($swaptotal) > 1 ? $swaptotal[1] : '';

// Extract the CPU model
$cpuinfo = file_get_contents("/proc/cpuinfo");
preg_match('/model name\s?:\s?(.+)$/m', $cpuinfo, $cpu_model);
$cpu_model = count($cpu_model) > 1 ? $cpu_model[1] : '';

echo json_encode([
  "hostname" => php_uname("n"),
  "user" => get_current_user(),
  "kernel" => php_uname("r"),
  "virtualization" => @exec("systemd-detect-virt"),
  "cloudlinux" => preg_match('/\.lve/', php_uname("r")) ? true : false,
  "io_schedulers" => $schedulers,
  "mounts" => file_get_contents("/proc/mounts"),
  "cpuinfo" => $cpuinfo,
  "cpu" => $cpu_model,
  "meminfo" => file_get_contents("/proc/meminfo"),
  "mem" => [
    "total" => $memtotal,
    "swap" => $swaptotal
  ],
  "users" => file_get_contents("/etc/passwd"),
  "phpinfo" => $phpinfo_capture,
  "db_version" => $mysql_version,
  "php_version" => phpversion(),
  "server" => array_key_exists('SERVER_SOFTWARE', $_SERVER) ? $_SERVER['SERVER_SOFTWARE'] : ''
]);

?>