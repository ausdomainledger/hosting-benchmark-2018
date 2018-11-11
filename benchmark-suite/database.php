<?php

// Begin configuration
$DB_DSN = "mysql:host=localhost;dbname=xxx_asic;charset=utf8mb4";
$DB_USER = "xxx_asic";
$DB_PASSWORD = "xxx_asic";
// End configuration

function get_data_if_needed() {
  if (file_exists("./asic.csv")) {
    return;
  }
  $fd = fopen("./asic.csv", "w");
  $ch = curl_init();
  curl_setopt_array($ch, [
    CURLOPT_FILE => $fd,
    CURLOPT_URL => "https://ausdomainledger.net/asic.csv",
    CURLOPT_TIMEOUT => 3600
  ]);
  if (!curl_exec($ch)) {
    die("Failed to download dataset");
  }
  curl_close($ch);
  fclose($fd);
  die("Downloaded dataset, run again");
}

function migrate_schema($db) {
  $migrations = [
    "DROP TABLE IF EXISTS asic_names;",
    <<<EOT
    CREATE TABLE asic_names (
      name VARCHAR(255) NOT NULL,
      status VARCHAR(255) NOT NULL,
      date_reg VARCHAR(255),
      date_cancelled VARCHAR(255),
      abn VARCHAR(255),
      FULLTEXT KEY idx_name (name)
    ) ENGINE=InnoDB;
EOT
  ];
  foreach ($migrations as $statement) {
    if (!$db->prepare($statement)->execute()) {
      die("$statement failed");
    }
  }
}

// Start
$db = new PDO($DB_DSN, $DB_USER, $DB_PASSWORD, [
  PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
  PDO::ATTR_EMULATE_PREPARES   => false,
  PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION
]);

migrate_schema($db);
get_data_if_needed();

// Insertion benchmark
$fd = fopen("./asic.csv", "r");
if (!$fd) {
  die("Can't read asic data");
}
fgets($fd, 1024); // Skip the header

// To be clear, this isn't a particularly fast method of insertion, 
// that's not what we're aiming for.
// We're going to do one statement per insert in order to keep this
// representative - bulk/batch inserts are not that typical and
// neither are multi-value statements.
$db->beginTransaction();
$stmt = $db->prepare("insert into asic_names values (?, ?, ?, ?, ?)");
$insert_row_count = 0;
$start = microtime(true);
while (($line = fgets($fd, 1024)) !== false) { // wc -L reports a much lower max line length on this dataset
  $data = explode("\t", $line);
  $stmt->execute([$data[1], $data[2], $data[3], $data[4], $data[8]]);
  $insert_row_count++;
  // We really don't need the full dataset, don't want to actually DoS the server.
  if ($insert_row_count >= 500e3) {
    break;
  }
}
fclose($fd);
$db->commit();
$insert_duration = microtime(true) - $start;

// Query Benchmark

// We need to generate random queries because e.g. Panthur run either a large query cache
// or a large caching MySQL proxy.
// For benchmarking this needs to be deterministic, so here we are:
mt_srand(1541667474505);

$stmt = $db->prepare("select * from asic_names where name LIKE ?");
$query_count = 0;
$start = microtime(true);
for ($i = 0; $i < PHP_INT_MAX; $i++) {
  if (microtime(true) - $start >= 15) {
    break;
  }

  $q = mt_rand(0, 1) === 1 ? '%' : chr(mt_rand(65, 90)); // Wildcard or capital
  for ($j = 0; $j < mt_rand(1, 5); $j++) {
    $q .= chr(mt_rand(97, 122));
  }
  $q .= mt_rand(0, 1) === 1 ? '%' : '';

  $stmt->execute([$q]);
  $stmt->fetchAll(PDO::FETCH_ASSOC);
  $query_count++;
}

$query_duration = microtime(true) - $start;

// Cleanup
$db = null;

echo json_encode([ "database.php" => [
  "insert_duration" => $insert_duration,
  "insert_rows" => $insert_row_count,
  "query_duration" => $query_duration,
  "queries" => $query_count
]], JSON_PRETTY_PRINT), PHP_EOL;

?>