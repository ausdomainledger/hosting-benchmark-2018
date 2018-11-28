<template>
<div id="app">
  <section class="main">
    <h1>Australian Shared Hosting Benchmark 2018</h1>
    <p>Published and retracted 2018-11-11, re-published 2018-11-28</p>

    <section class="intro">
      <h2 id="introduction">Introduction</h2>
      <p>The purpose of this project was to figure out the following:</p>
      <blockquote>
        Who, of the most popular<a id="r1" href="#fn1">[1]</a> Australian web hosts,
        provides the fastest cPanel hosting at the $20/month<a id="r2" href="#fn2">[2]</a>
        price point?
      </blockquote>
      <p>It follows up on a similar question I asked in 2017:
        <a href="https://ausdomainledger.net/cpanelhosts2017/" target="_blank">
          Who is the biggest cPanel web host in Australia?
        </a>
      </p>
      <p>
        It was also important for the benchmark to clearly justify all of its decisions and
        to provide a way for others to reproduce the results on the same or similar
        kinds of services. To that end,
        <a href="https://github.com/ausdomainledger/hosting-benchmark-2018" target="_blank"
          rel="noopener noreferrer">
        every part of this benchmark is open source</a>,
        including the tools and raw data, so you can run any or all parts of it yourself.
      </p>
      <p><strong>A caution:</strong> Remember, these results only reflect around a week
      (~4 daily runs per provider) of testing,
      Providers update their services constantly, and all benchmarks in existence are flawed.
      Consider that when choosing how seriously to take these results in your decision-making.</p>
    </section>
    <section class="services">
      <h2 id="services">Providers and Services</h2>
      <p>The specific cPanel services that were benchmarked were selected according to the
      <a href="#methodology">Methodology</a>, which can be approximately be summed up as:
      take the ~10 most popular web hosts on Whirlpool and choose their ~$20/mo plan.</p>
      <p>These are then benchmarked against a VPS service at the same price,
        configured in a performance-optimal way<a id="r3" href="#fn3">[3]</a>.</p>
      <div class="scrolling-table">
        <table>
          <thead>
            <tr>
              <th @click="toggleProviderSort(p => p.plan.provider)">Provider</th>
              <th @click="toggleProviderSort(p => p.plan.name)" style="min-width: 150px;">
                Plan Name</th>
              <th @click="toggleProviderSort(p => p.plan.cost)">Cost/Month</th>
              <th @click="toggleProviderSort(p => p.env.hostname)">Hostname</th>
              <th @click="toggleProviderSort(p => p.env.cloudlinux)">CloudLinux</th>
              <th @click="toggleProviderSort(p => p.env.server)">Webserver</th>
              <th @click="toggleProviderSort(p => p.env.php_version)">PHP Version</th>
              <th @click="toggleProviderSort(p => p.env.db_version)">Database Version</th>
              <th @click="toggleProviderSort(p => p.env.kernel)">Kernel</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="provider in sortedProviders" v-bind:key="provider.provider"
              v-bind:class="provider.provider">
              <td class="logo" v-bind:class="provider.provider">
                {{ provider.plan.provider }}
                <div v-if="provider === baseline">
                  <small><strong>(VPS Benchmark)</strong></small>
                </div>
              </td>
              <td>{{ provider.plan.name }}</td>
              <td>${{ (provider.plan.cost / 100).toFixed(2) }}</td>
              <td style="word-break: break-all;">{{ provider.env.hostname }}</td>
              <td>{{ provider.env.cloudlinux === true ? 'Yes' : 'No' }}</td>
              <td>{{ provider.env.server }}</td>
              <td>{{ provider.env.php_version }}</td>
              <td>{{ provider.env.db_version }}</td>
              <td>{{ provider.env.kernel }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p>
        Synergy Wholesale, TPP Wholesale, Crazy Domains, and Net Registry were
          all disqualified<a href="#fn7" id="r7">[7]</a>, bringing the field down
          to 8 providers.
      </p>
    </section>
    <section class="winner">
      <h2 id="winner">Overall Winner: Netorigin</h2>
      <img src="./assets/netorigin_logo.svg" alt="Net Virtue" style="height: 50px;">
      <p>The overall winner is the Perth-based Netorigin's "Elite" cPanel service.
         It performed outstandingly on all benchmarks, exceeding the VPS in a number of
         cases.
         &#127881;
      </p>
      <h3>and on the eastern coast ...</h3>
      <img src="./assets/ventraip_logo.svg" alt="VentraIP" style="height: 50px;">
      <p>Since, as we all know, nobody lives in Perth, it wouldn't be fun unless we also
        crowned an east-coast winner.
        That glory belongs to VentraIP, who came close enough to Netorigin on the WP
        Overall Score, that you could argue for a draw.
      </p>
    </section>
    <section class="results">
      <h2 id="results">Results</h2>
      <div class="scrolling-table">
        <table>
          <thead>
            <tr>
              <th @click="toggleResultsSort(p => p.plan.provider)">Provider</th>
              <th style="min-width: 100px;"
              @click="toggleResultsSort(p => wpOverall(p))">
                WP Overall Score</th>
              <th style="min-width: 100px;"
              @click="toggleResultsSort(p => p.measurements.WP1P90.Median)">
                WP P90 1/sec</th>
              <th style="min-width: 100px;"
              @click="toggleResultsSort(p => p.measurements.WP5P90.Median)">
                WP P90 5/sec</th>
              <th @click="toggleResultsSort(p => p.measurements.WP10P90.Median)">
                WP P90 10/sec</th>
              <th @click="toggleResultsSort(p => p.measurements.DBInsertDuration.Median)">
                DB Insert</th>
              <th @click="toggleResultsSort(p => p.measurements.DBQueries.Median)">
                DB Queries</th>
              <th @click="toggleResultsSort(p => p.measurements.CPUOps.Median)">
                Prime Sieve</th>
              <th @click="toggleResultsSort(p => p.measurements.PHP1Duration.Median)">
                PHP.net #1</th>
              <th @click="toggleResultsSort(p => p.measurements.PHP2Duration.Median)">
                PHP.net #2</th>
              <th @click="toggleResultsSort(p => p.measurements.IORandomRW.Median)">
                IO Random</th>
              <th @click="toggleResultsSort(p => p.measurements.IOOpenDuration.Median)">
                IO Open</th>
              <th @click="toggleResultsSort(p => p.measurements.IOSeqWrite.Median)">

                IO Seq. Write</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="provider in sortedResults" v-bind:key="provider.provider"
              v-bind:class="provider.provider">
              <td class="logo" v-bind:class="provider.provider">
                {{ provider.plan.provider }}
                <div v-if="provider === baseline">
                  <small><strong>(VPS Benchmark)</strong></small>
                </div>
              </td>
              <td>
                <span v-if="provider === baseline">100%</span>
                <div v-else>
                  {{ diff(wpOverall(provider), wpOverall(baseline)) }}
                  overall
                </div>
              </td>
              <td>
                {{ (provider.measurements.WP1P90.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.WP1P90.Median,
                  baseline.measurements.WP1P90.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.WP5P90.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.WP5P90.Median,
                  baseline.measurements.WP5P90.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.WP10P90.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.WP10P90.Median,
                  baseline.measurements.WP10P90.Median) }}
                </small>
              </td>
              <td>
                {{ ((provider.measurements.DBInsertDuration.Median / 1e6) / 1000).toFixed(2) }}s
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.DBInsertDuration.Median,
                  baseline.measurements.DBInsertDuration.Median) }}
                </small>
              </td>
              <td>
                {{ provider.measurements.DBQueries.Median }}
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(baseline.measurements.DBQueries.Median,
                  provider.measurements.DBQueries.Median) }}
                </small>
              </td>
              <td>
                {{ provider.measurements.CPUOps.Median }}
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(baseline.measurements.CPUOps.Median,
                  provider.measurements.CPUOps.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.PHP1Duration.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.PHP1Duration.Median,
                  baseline.measurements.PHP1Duration.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.PHP2Duration.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.PHP2Duration.Median,
                  baseline.measurements.PHP2Duration.Median) }}
                </small>
              </td>
              <td>
                {{ provider.measurements.IORandomRW.Median }}
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(baseline.measurements.IORandomRW.Median,
                  provider.measurements.IORandomRW.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.IOOpenDuration.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.IOOpenDuration.Median,
                  baseline.measurements.IOOpenDuration.Median) }}
                </small>
              </td>
              <td>
                {{ (provider.measurements.IOSeqWrite.Median / 1e6).toFixed(2) }}ms
                <small style="display: block;" v-if="provider !== baseline">
                  {{ diff(provider.measurements.IOSeqWrite.Median,
                  baseline.measurements.IOSeqWrite.Median) }}
                </small>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <h3>Visualized (based on WP Overall Score)</h3>
      <div class="barcharts">
        <div v-for="p in sortedResults" v-bind:key="p.provider" class="provider">
          <div class="bar" :style="'width:'
          + ((wpOverall(baseline) / wpOverall(p)) * 100) + '%;'"></div>
          <div class="name" v-if="p != baseline">{{ p.plan.provider }}</div>
        </div>
      </div>
    </section>
    <section class="methodology">
      <h2 id="methodology">Methodology</h2>
      <h3>Provider Selection</h3>
      <p>The selection of providers is a task that could be strongly
        affected by bias, so I endeavored to give it some semblance of being
        data-driven.</p>
      <p>The ideal scenario was to scrape all of the threads posted in during 2018 in the
        <a href="http://forums.whirlpool.net.au/forum/116" target="_blank">
        Whirlpool Web Hosting forum</a>, apply named-entity recognition on the
        posts and then perform a frequency count for each named-entity (web host).
      </p>
      <p>Unfortunately I did not quite reach that scenario. While the scraping of
        the forum was easy to perform (and with permission), the state-of-the-art
        for NER requires a significant amount of training to be at all effective.
        e.g. "net registry" is very hard to pick up as an NE. In the end, I
        (very unglamorously) used my eyeballs to identify the top 17 web hosts from
        the ~16k posts.
      </p>
      <p>From there, permissive RegExps were written for those 17
        and used to generate the frequencies.</p>
      <p>One significant mistake that I made was counting all posts whether they were
        made in 2018 or not (so threads that go for many years had all posts counted).
        A second mistake was that posts from another forum ("Server Management")
        accidentally leaked into the dataset during the scrape, but I suspect
        its impact was negligible.
      </p>
      <p>You can look in the source code repository for all of the raw data and the
        actual frequencies.
      </p>
      <h3>Benchmark Selection</h3>
      <p>A meaningful measure of "performance" in shared web hosting should ideally
        reflect the speed at which humans experience the web applications that
        they run on their cPanel service. As such, we will run only pure-PHP benchmarks
        which are initiated over HTTP requests. This rules out usual tooling like
        sysbench, FIO, etc. Another reason for using pure PHP is that a number of hosts
        have disabled system/exec and it's not clear whether it's possible to fork
        an arbitrary ELF executable (e.g. using cgi-bin) across every host or not. Since it
        it has a chance of torpedoing the entire methodology, I chose to avoid the risk.
      </p>
      <p>In benchmarking, it's unavoidable to rely on microbenchmarks. Even for
        "realistic" tests, they will not truthfully predict a real-world workload.
        However, this is a comparative study. We are not gunning for accuracy
        in the absolute numbers reported by the benchmarks, just for accuracy in
        the comparative numbers. For this reason,
        the results are shown as a percentage of the
        performance of a baseline environment (a fast VPS).</p>
      <p>The baseline environment is an unmanaged Linux VPS from Binary Lane,
        at the same cost target as the shared hosting plans. It is configured
        in a performance-maximizing way (nginx+PHP-FPM), with the aim of providing
        some guidance as to "how much slower/faster" a shared service is, compared
        to an optimum (though not necessarily <em>the</em> optimum
        <a id="r6" href="#fn6">[6]</a>).
      </p>
      <p>Each benchmark is also repeated 25-30 across a week time interval, in order to account
        for transient "noisy neighbor" effects.
      </p>
      <p>The pure PHP benchmarks are described below.</p>
      <h4>1. WordPress</h4>
      <p>This is the "realistic" benchmark that primarily informs the overall result.</p>
      <p>It involves installing an entirely vanilla<a id="r4" href="#fn4">[4]</a>
        WordPress installation
        and performing a constant-rate, concurrent HTTP siege against it using wrk2, the
        "coordinated-omission" fork of wrk. I have
        <a href="https://github.com/ausdomainledger/wrk2/compare/master...ausdomainledger:multi-ip"
        target="_blank" rel="noopener noreferrer">further forked wrk2 at the
        ausdomainledger/wrk2 multi-ip branch</a> in order to round-robin requests between
        IPv4 and IPv6 addresses, when available. The reason for this is that Zuver and
        VentraIP have extremely aggressive CSF firewall rules configured, which causes
        the 10/sec concurrency level to fail. Using multiple source IPs permits us to
        evade the firewall.
      </p>
      <p>The aim of the benchmark is to record the latency percentiles for the site's
        transaction time at fixed concurrency levels (1/sec, 5/sec, 10/sec). These should
        approximately reflect the usability of the website when it is being browsed
        by only the webmaster, by a few users, and a few more users, respectively.
        Initially 10/sec was 25/sec, but it became clear during benchmarking that all
        but one of the providers fell flat on their faces under such a workload.
      </p>
      <p>Due to shared hosting being a multi-tenant environment, we will focus mainly
        on the P90 result, as anything higher is more or less just noise. However,
        you may browse the raw data for a complete latency distribution for each
        provider.
      </p>
      <p>The "siege" is performed from a dedicated server hosted on Servers Australia's Sydney
        network,
        and is within a millisecond latency of every benchmarked website. The first exception is
        Micron21, who are hosted in Melbourne, who are sieged from a separate VPS node hosted
        in Micron21's Kilsyth DC. The second exception is Net Origin, who are hosted in Perth,
        who are sieged from a separate VPS node hosted by RansomIT in Perth. In both cases,
        extra care was taken to ensure that the different environment did not introduce
        extra advantage or penalty.
      </p>
      <h4>2. Database</h4>
      <p>This is a second "realistic" benchmark which focuses on the ability
        of the user to put the database server to work. It is significant to include
        because each host may or may not be running
        <a href="https://docs.cloudlinux.com/mysql_governor.html" target="_blank">
        CloudLinux's MySQL governor</a>.
      </p>
      <p>The benchmark is based upon the first 500k records in the ASIC Business
        Names dataset, with a full-text MySQL search index on the business name column.
        The benchmark itself reports on the throughput of INSERTs and full-text queries.
      </p>
      <h4>3. Microbenchmark: CPU - Prime Sieve</h4>
      <p>This is a better CPU microbenchmark than #4 and #5, basically because it runs
        longer. It is my own port of Sieve of Eratosthenes used
        by sysbench to measure single-threaded CPU performance.
      </p>
      <h4>4. Microbenchmark: PHP.net's bench.php</h4>
      <p>This is one of the two official microbenchmarks used by PHP to report on
        performance gains and regressions between PHP versions. However, it runs very
        quickly and is not very good at measuring any meaningful workload.
      </p>
      <h4>5. Microbenchmark: PHP.net's micro_bench.php</h4>
      <p>As above.</p>
      <h4>6. Microbenchmark: Disk IO</h4>
      <p>Disk IO is a terribly important subsystem when measuring performance, but
      it is more or less impossible<a href="#fn5" id="r5">[5]</a>
      to accurately measure from a pure PHP benchmark. For this reason, it is heavily
      discounted in the overall result.
      </p>
      <p>Nonetheless, an attempt is given to measure this, for curiosity, if no
        other reason. The benchmark essentially measures mass random read/write/open
        on inodes characteristic of the kind you would find in a shared hosting
        web application (e.g. WordPress).
      </p>
      <h4>bUt iS iT sTaTisTiCaLlY sIgNiFiCanT?</h4>
      <p>Directly comparing samples (i.e. percentiles) across datasets is potentially
        problematic. I haven't proven that the measured differences did not occur
        purely by chance.
      </p>
      <p>By running the benchmark suites 25-30 times across different days and
        for significant durations, I have tried to provide an accurate basis for
        comparison, but my good feeling about it is purely intuitive rather than
        benefiting from rigorous stats analysis.
      </p>
      <h4>Ranking the Results</h4>
      <p>
        For the overall rankings, the sum of the median of each of 1/sec, 5/sec and 10/sec
        P90s is taken.
      </p>
      <p>
        It might be surprising to ignore the other measurements entirely,
        but coming up with a
        composite/summary score based on some kind of factor or
        principal component analysis proved to
        be too unwieldly for me to do correctly, and in a reasonable time.
      </p>
      <p>
        This ranking method is also the most direct measure for the idea of "performance"
        that most people think about in the context of shared hosting, in an environment
        that is more or less homogenous towards and optimized for WordPress workloads.
      </p>
      <p>
        In any case, you are free to sort the results by whatever measurement you want.
      </p>
    </section>
    <section class="disclosures">
      <h2 id="disclosures">Disclosures and Contact</h2>
      <p>Of all of these providers, I am affiliated only with the baseline
        benchmark provider, Binary Lane, as an existing customer. Their
        selection as the benchmark provider was made on the basis of personal
        preference and the fact that they are not eligible as a shared
        hosting provider.
      </p>
      <p>No providers were informed about being benchmarked at any time.
        Care was taken not to reveal the purpose of the accounts.
      </p>
      <p>All providers were paid the full amount for their services using only
        my personal funds.
      </p>
      <p>Services were used only for the purpose of this benchmark, and
        then cancelled. I received an automatic refund from VentraIP &amp;
        Zuver at cancellation time, apparently in accordance with their
        45 day money-back guarantee.
      </p>
      <p>Feedback can be provided by &#128231; to _&lt;at&gt;ausdomainledger.net.
        I am happy to post factual corrections or highlight mistakes and flaws. However,
        no existing or additional benchmarks will be re-run or have their results
        updated.
      </p>
      <p>All product names, logos, and brands are property of their respective owners.
        All company, product and service names used in this website are for identification
        purposes only. Use of these names, logos, and brands does not imply endorsement.</p>
    </section>
    <section class="footnotes">
      <hr>
      <p id="fn1">
        <a href="#r1">[1]</a> Popularity is calculated by the number of mentions on the
        Whirlpool Web Hosting forum in 2018. See Methodology for more info.
      </p>
      <p id="fn2">
        <a href="#r2">[2]</a> $20/month was selected as an arbitrary, comfortable price-point,
        that feels like the mean for what is offered by hosts and chosen by customers,
        but isn't backed by data.
      </p>
      <p id="fn3">
        <a href="#r3">[3]</a> The Ansible provisioning script for the VPS baseline benchmark
        can be found <a href="https://github.com/ausdomainledger/hosting-benchmark-2018"
        target="_blank" rel="noopener noreferrer">
        in the source code repository</a>.
      </p>
      <p id="fn4">
        <a href="#r4">[4]</a> Caching plugins are entirely orthogonal to the purpose of this
        benchmark. We are trying to perform a comparative study of the "raw power" of each
        cPanel service. To that end, any caching plugins that are snuck in by the server
        (like Litespeed Cache) are forcibly removed.
      </p>
      <p id="fn5">
        <a href="#r5">[5]</a> PHP does not expose a sufficient set of syscalls to make it
        possible to control whether a write is actually applied to disk. At most, disk writes
        will result in dirty pages in the kernel, so the benchmark results will ultimately
        depend on the state of the kernel page cache at any point in time (in other words,
        completely out of our control).
      </p>
      <p id="fn6">
        <a href="#r6">[6]</a> The selection of Binary Lane as the benchmark baseline is not
        any kind of endorsement or statement about its worthiness as a web host. It is purely
        a question of convenience and independence from the benchmarked providers.
      </p>
      <p id="fn7">
        <a href="#r7">[7]</a> Synergy Wholesale and TPP Wholesale were disqualified for not
        having any retail (not requiring pre-approval/contract signing) services available.
        Crazy Domains were disqualified for requiring annual+ payment terms, which does not fit
        the spirit of the problem statement. Net Registry were disqualified for
        unironically asking me to fax them my drivers licence, but it is unlikely that their
        results would differ from that of Melbourne IT by much.
      </p>
    </section>
  </section>
</div>
</template>

<script>
import benchmarkResults from '@/results.json';

const rankingFunctions = {
  'wp-median': p => p.measurements.WP1P90.Median
    + p.measurements.WP5P90.Median
    + p.measurements.WP10P90.Median,
};

export default {
  name: 'app',
  components: {
  },
  data() {
    return {
      providerSortKey: p => p.plan.provider,
      providerSortAsc: true,
      sortKey: rankingFunctions['wp-median'],
      sortDir: true,
      results: Object.values(benchmarkResults),
      wpOverall: rankingFunctions['wp-median'],
    };
  },
  created() {
  },
  methods: {
    toggleProviderSort(accessorFn) {
      if (accessorFn(this.results[0]) === this.providerSortKey(this.results[0])) {
        this.providerSortAsc = !this.providerSortAsc;
      } else {
        this.providerSortKey = accessorFn;
        this.providerSortAsc = true;
      }
    },
    toggleResultsSort(accessorFn) {
      if (accessorFn(this.results[0]) === this.sortKey(this.results[0])) {
        this.sortDir = !this.sortDir;
      } else {
        this.sortKey = accessorFn;
        this.sortDir = true;
      }
    },
    diff(a, b) {
      const d = a / b;
      if (d > 1) {
        return `${d.toFixed(2)}x slower`;
      }
      if (d < 1) {
        return `${(b / a).toFixed(2)}x faster`;
      }
      return 'the same';
    },
  },
  computed: {
    baseline() {
      return this.results[0];
    },
    sortedResults() {
      return this.results.slice().sort((a, b) => {
        const aVal = this.sortKey(a);
        const bVal = this.sortKey(b);
        if (typeof aVal === 'string') {
          return this.sortDir
            ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
        }
        return this.sortDir ? aVal - bVal : bVal - aVal;
      });
    },
    sortedProviders() {
      return this.results.slice().sort((a, b) => {
        const aVal = this.providerSortKey(a);
        const bVal = this.providerSortKey(b);
        if (typeof aVal === 'string') {
          return this.providerSortAsc
            ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
        }
        return this.providerSortAsc ? aVal - bVal : bVal - aVal;
      });
    },
  },
};
</script>

<style lang="scss">
html {
  font-size: 18px;
  line-height: 30px;
  scroll-behavior: smooth;
}
body {
  margin: 0;
  padding: 0;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
  max-width: 1400px;
  margin: 0 auto;
}
h2, h3 {
  border-bottom: 1px dotted #2c3e50;
}
h1 {
  font-size: 2rem;
}
h2 {
  font-size: 1.5rem;
}
h3 {
  font-size: 1.2rem;
}
h4 {
  font-size: 1.1rem;
  margin: 0 auto;
}
p {
  margin-bottom: 1rem;
}
section {
  margin: 2.5rem auto;
}
.main {
  margin: auto;
}
.footnotes {
  hr {
    border: 1px dashed lighten(#2c3e50, 50);
  }
  font-size: 0.9rem;
}
blockquote {
  font-size: 1.1rem;
}
span[title] {
  border-bottom: 1px #2c3e50 dashed;
}
.scrolling-table {
  overflow-x: scroll;
  table {
    width: 100%;
    font-size: 0.9rem;
    border-collapse: collapse;
    thead {
      background: #2c3e50;
      color: white;
      cursor: pointer;
      user-select: none;
    }
    tbody tr:nth-child(odd) {
      background: #eee;
    }
    tr.binarylane {
      font-style: italic;
      border-bottom: 1px dotted #2c3e50 !important;
    }
    td, th {
      text-align: left;
      line-height: initial;
      padding: 0 1rem;
      &.logo {
        min-width: 100px;
      }
      &.logo img {
        max-width: 150px;
        max-height: 50px;
      }
    }
  }
}
.benchmark-row {
  background: whitesmoke;
  font-style: italic;
}
.winner {
  background-color: #2c3e50;
  color: white;
  padding: 1em;
  h2 {
    border-bottom: 1px dotted white;
  }
}
.barcharts {
  margin-top: 1rem;
  .provider {
    &:nth-child(1) {
      .bar {
        background: #2c3e50;
        padding-right: 1rem;
        text-align: right;
        &::after {
          content: "VPS Benchmark";
          color: white;
        }
      }
    }
    margin-bottom: 0.25rem;
    display: flex;
    flex-direction: row;
    .bar {
      background: lighten(#2c3e50, 50%);
    }
    .name {
      margin-left: 1rem;
    }
  }
}
</style>
