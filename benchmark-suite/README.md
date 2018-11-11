# Benchmarks

| Name | Description | Authorship | 
|------|-------------|------------|
WordPress `wordpress.sh` | Multi-threaded "realistic" benchmark focused on latency percentiles under a variety of fixed visitor loads, for a completely vanilla WordPress site. | Own
Database `database.php` | Single-threaded PHP + MySQL benchmark focused on a database-intensive workload using 500k records from the [ASIC Business Names dataset](https://data.gov.au/dataset/asic-business-names). | Own
CPU - PHP.net's `bench.php` | Unmodified short-duration single-threaded official PHP microbenchmark focused on calculation speed (https://github.com/php/php-src/blob/php-7.2.11/Zend/bench.php)  | PHP Official |
CPU - PHP.net's `micro_bench.php` | Unmodified short-duration single-threaded official PHP microbenchmark focused on calculation speed (https://github.com/php/php-src/blob/php-7.2.11/Zend/micro_bench.php)  | PHP Official | 
CPU - Prime Sieve `micro_cpu.php` | Longer-duration single-threaded PHP microbenchmark using sieve of eratosthenes | Own, [sieve ported from sysbench CPU test](https://github.com/akopytov/sysbench/blob/master/src/tests/cpu/sb_cpu.c#L109) |
Disk IO `micro_io.php` | Single-threaded PHP microbenchmark focused on IO. Since PHP does not expose the right set of syscalls to be able sufficiently control the kernel page cache and physical disk writes, this benchmark is of questionable value. | Own