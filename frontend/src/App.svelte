<script>
  import { onMount } from 'svelte';
  import { api } from './api.js';
  import IPTimeline from './components/IPTimeline.svelte';

  let vpcCidrInput = '';
  let subnetCidrInput = '';
  let loading = false;
  let error = '';
  let analysisResult = null;
  let abortController = null; // Track current request for cancellation

  // Example CIDR ranges for demonstration
  const exampleVpcCIDRs = [
    '10.0.0.0/16',
    '172.16.0.0/20'
  ].join('\n');

  const exampleSubnetCIDRs = [
    '10.0.1.0/24',
    '10.0.2.0/24',
    '172.16.1.0/24',
    '172.16.3.0/24'
  ].join('\n');

  // Set example CIDRs on component mount
  onMount(() => {
    vpcCidrInput = exampleVpcCIDRs;
    subnetCidrInput = exampleSubnetCIDRs;
  });

  // Parse CIDR inputs into arrays
  function parseVpcCIDRInput() {
    console.log('parseVpcCIDRInput called with:', vpcCidrInput);
    const result = vpcCidrInput
      .split('\n')
      .map(line => line.trim())
      .filter(line => line.length > 0);
    console.log('parseVpcCIDRInput result:', result);
    return result;
  }

  function parseSubnetCIDRInput() {
    console.log('parseSubnetCIDRInput called with:', subnetCidrInput);
    const result = subnetCidrInput
      .split('\n')
      .map(line => line.trim())
      .filter(line => line.length > 0);
    console.log('parseSubnetCIDRInput result:', result);
    return result;
  }

  // Analyze CIDR ranges
  async function analyzeCIDRs() {
    // Cancel any previous request
    if (abortController) {
      abortController.abort();
    }

    // Small delay to ensure input values are fully updated
    await new Promise(resolve => setTimeout(resolve, 10));

    // Debug: Log current input values
    console.log('=== ANALYSIS START ===');
    console.log('Raw VPC Input:', JSON.stringify(vpcCidrInput));
    console.log('Raw Subnet Input:', JSON.stringify(subnetCidrInput));

    const vpcCidrs = parseVpcCIDRInput();
    const subnetCidrs = parseSubnetCIDRInput();
    const allCidrs = [...vpcCidrs, ...subnetCidrs];
    
    // Debug: Log parsed values
    console.log('Parsed VPC CIDRs:', JSON.stringify(vpcCidrs));
    console.log('Parsed Subnet CIDRs:', JSON.stringify(subnetCidrs));
    console.log('All CIDRs Combined:', JSON.stringify(allCidrs));
    
    if (allCidrs.length === 0) {
      error = 'Please enter at least one CIDR range';
      return;
    }

    // Prevent multiple simultaneous requests
    if (loading) {
      return;
    }

    loading = true;
    error = '';
    analysisResult = null; // Clear previous results immediately
    
    try {
      console.log('Sending VPC CIDRs:', vpcCidrs);
      console.log('Sending Subnet CIDRs:', subnetCidrs);
      console.log('Sending All CIDRs to API:', allCidrs);
      
      const requestBody = { 
        cidrs: allCidrs,
        vpc_cidrs: vpcCidrs,
        subnet_cidrs: subnetCidrs,
        timestamp: Date.now() // Add timestamp to prevent caching
      };
      
      console.log('Request body:', JSON.stringify(requestBody, null, 2));
      
      // Create abort controller for cancellation
      abortController = new AbortController();
      
      // Use the api module to send the request
      analysisResult = await api.analyzeCIDRs(requestBody, { 
        signal: abortController.signal 
      });
      
      // Debug: Log received analysis result before modification
      console.log('Received analysis result from API:', JSON.stringify(analysisResult, null, 2));
      
      // Add the categorization to the result for the timeline component
      analysisResult.vpc_cidrs = vpcCidrs;
      analysisResult.subnet_cidrs = subnetCidrs;
      
      console.log('Final analysis result with categories added:', analysisResult);
      console.log('=== ANALYSIS COMPLETE ===');
    } catch (err) {
      console.error('Analysis error:', err);
      if (err.name === 'AbortError') {
        error = 'Analysis timed out. Please try again or reduce the number of CIDR ranges.';
      } else if (err.message.includes('Failed to fetch')) {
        error = 'Failed to connect to the analysis service. Please ensure the backend is running.';
      } else {
        error = `Analysis failed: ${err.message}`;
      }
      analysisResult = null;
    } finally {
      loading = false;
      abortController = null;
    }
  }

  // Clear all data
  function clearAll() {
    // Cancel any ongoing request
    if (abortController) {
      abortController.abort();
    }
    
    vpcCidrInput = '';
    subnetCidrInput = '';
    analysisResult = null;
    error = '';
    loading = false;
  }

  // Cancel the current analysis
  function cancelAnalysis() {
    if (abortController) {
      abortController.abort();
    }
    loading = false;
    error = 'Analysis cancelled by user';
  }

  // Load example CIDRs
  function loadExample() {
    vpcCidrInput = exampleVpcCIDRs;
    subnetCidrInput = exampleSubnetCIDRs;
    analysisResult = null;
    error = '';
  }

  // Load user's test data from test.dat
  function loadTestData() {
    // Based on the user's test.dat file
    vpcCidrInput = `10.59.204.128/25
10.58.42.128/26`;
    
    subnetCidrInput = `10.59.204.128/26
10.59.204.192/26
10.58.204.64/26
10.58.204.96/27
10.58.204.64/27
10.58.42.160/27
10.58.42.128/27`;
    
    analysisResult = null;
    error = '';
    console.log('Loaded test data from test.dat');
  }
</script>

<main>
  <header>
    <h1>🌐 CIDR Viewer</h1>
    <p>Visualize IPv4 CIDR ranges and identify gaps and unused IP addresses</p>
  </header>

  <div class="container">
    <!-- Input Panel -->
    <div class="input-panel">
      <h2>Network Configuration</h2>
      
      <div class="input-section">
        <h3>VPC CIDR Blocks</h3>
        <p>Enter VPC-level CIDR ranges (one per line):</p>
        <textarea
          class="textarea"
          bind:value={vpcCidrInput}
          placeholder="Enter VPC CIDR ranges, e.g.:&#10;10.0.0.0/16&#10;172.16.0.0/20"
        ></textarea>
      </div>

      <div class="input-section">
        <h3>Subnet CIDRs</h3>
        <p>Enter subnet CIDR ranges within VPCs (one per line):</p>
        <textarea
          class="textarea"
          bind:value={subnetCidrInput}
          placeholder="Enter subnet CIDR ranges, e.g.:&#10;10.0.1.0/24&#10;10.0.2.0/24&#10;172.16.1.0/24"
        ></textarea>
      </div>

      <div>
        {#if loading}
          <button class="button" on:click={cancelAnalysis}>Cancel Analysis</button>
        {:else}
          <button class="button" on:click={analyzeCIDRs} disabled={loading}>
            Analyze Network
          </button>
          <button class="button" on:click={loadExample}>Load Example</button>
          <button class="button" on:click={loadTestData}>Load Test Data</button>
          <button class="button" on:click={clearAll}>Clear All</button>
        {/if}
      </div>

      {#if error}
        <div class="error-message">{error}</div>
      {/if}
    </div>

    <!-- Visualization Panel -->
    <div class="visualization-panel">
      <h2>Analysis Results</h2>
      
      {#if loading}
        <div class="loading">
          <div class="spinner"></div>
          <div class="loading-text">
            <p><strong>Analyzing CIDR ranges...</strong></p>
            <p>Processing network topology and identifying gaps/overlaps</p>
            <p><small>This may take a few moments for large networks</small></p>
          </div>
        </div>
      {:else if analysisResult}
        <!-- Summary -->
        {#if analysisResult.summary}
        <div>
          <h3>Summary</h3>
          <table class="table">
            <tr>
              <td><strong>Total IPs:</strong></td>
              <td>{analysisResult.summary.total_ips.toLocaleString()}</td>
            </tr>
            <tr>
              <td><strong>Allocated IPs:</strong></td>
              <td class="status-valid">{analysisResult.summary.allocated_ips.toLocaleString()}</td>
            </tr>
            <tr>
              <td><strong>Available IPs (Gaps):</strong></td>
              <td class="status-gap">{analysisResult.summary.available_ips.toLocaleString()}</td>
            </tr>
            <tr>
              <td><strong>Gap Count:</strong></td>
              <td class="status-gap">{analysisResult.summary.gap_count}</td>
            </tr>
            <tr>
              <td><strong>Overlap Count:</strong></td>
              <td class="status-overlap">{analysisResult.summary.overlap_count}</td>
            </tr>
          </table>
        </div>
        {/if}

        <!-- IP Timeline Visualization -->
        <IPTimeline {analysisResult} />

        <!-- Valid CIDRs -->
        {#if analysisResult.valid_cidrs && analysisResult.valid_cidrs.length > 0}
          <div>
            <h3>Valid CIDR Ranges</h3>
            <table class="table">
              <thead>
                <tr>
                  <th>CIDR</th>
                  <th>Network</th>
                  <th>Broadcast</th>
                  <th>Total IPs</th>
                  <th>Usable IPs</th>
                </tr>
              </thead>
              <tbody>
                {#each analysisResult.valid_cidrs as cidr}
                  <tr>
                    <td class="status-valid">{cidr.original}</td>
                    <td>{cidr.network}</td>
                    <td>{cidr.broadcast}</td>
                    <td>{cidr.total_ips.toLocaleString()}</td>
                    <td>{cidr.usable_ips.toLocaleString()}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        <!-- Invalid CIDRs -->
        {#if analysisResult.invalid_cidrs && analysisResult.invalid_cidrs.length > 0}
          <div>
            <h3>Invalid CIDR Ranges</h3>
            <table class="table">
              <thead>
                <tr>
                  <th>CIDR</th>
                  <th>Error</th>
                </tr>
              </thead>
              <tbody>
                {#each analysisResult.invalid_cidrs as cidr}
                  <tr>
                    <td class="status-invalid">{cidr.original}</td>
                    <td class="status-invalid">{cidr.error_msg}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        <!-- Gaps -->
        {#if analysisResult.gaps && analysisResult.gaps.length > 0}
          <div>
            <h3>IP Gaps (Unused Ranges)</h3>
            <table class="table">
              <thead>
                <tr>
                  <th>Start IP</th>
                  <th>End IP</th>
                  <th>Size</th>
                  <th>Suggested CIDR</th>
                </tr>
              </thead>
              <tbody>
                {#each analysisResult.gaps as gap}
                  <tr>
                    <td class="status-gap">{gap.start_ip}</td>
                    <td class="status-gap">{gap.end_ip}</td>
                    <td>{gap.size.toLocaleString()}</td>
                    <td>{gap.suggested_cidr}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        <!-- Overlaps -->
        {#if analysisResult.overlaps && analysisResult.overlaps.length > 0}
          <div>
            <h3>Overlapping Ranges</h3>
            <table class="table">
              <thead>
                <tr>
                  <th>CIDR 1</th>
                  <th>CIDR 2</th>
                  <th>Type</th>
                  <th>Intersection</th>
                </tr>
              </thead>
              <tbody>
                {#each analysisResult.overlaps as overlap}
                  <tr>
                    <td class="status-overlap">{overlap.cidr1}</td>
                    <td class="status-overlap">{overlap.cidr2}</td>
                    <td>{overlap.type}</td>
                    <td>{overlap.intersection}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        {#if analysisResult && analysisResult.gaps && analysisResult.overlaps && analysisResult.gaps.length === 0 && analysisResult.overlaps.length === 0}
          <div class="status-valid">
            <h3>✅ Perfect! No gaps or overlaps detected.</h3>
          </div>
        {/if}
      {:else}
        <div>
          <p>Enter CIDR ranges in the input panel and click "Analyze CIDRs" to see the results.</p>
          <p><strong>Example CIDR ranges:</strong></p>
          <ul>
            <li>192.168.1.0/24 (256 IPs)</li>
            <li>10.0.0.0/16 (65,536 IPs)</li>
            <li>172.16.0.0/20 (4,096 IPs)</li>
          </ul>
        </div>
      {/if}
    </div>
  </div>
</main>

<style>
  main {
    text-align: left;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
  }

  header h1 {
    font-size: 2.5rem;
    margin: 0;
    color: #646cff;
  }

  header p {
    font-size: 1.1rem;
    color: #888;
    margin: 0.5rem 0 0 0;
  }

  h2 {
    color: #646cff;
    margin-top: 0;
  }

  h3 {
    color: #888;
    margin-top: 1.5rem;
    margin-bottom: 0.5rem;
  }

  .input-section {
    margin-bottom: 1.5rem;
  }

  .input-section h3 {
    margin-top: 0;
    font-size: 1.1rem;
    color: #646cff;
  }

  .input-section p {
    margin: 0.25rem 0 0.5rem 0;
    font-size: 0.875rem;
    color: #888;
  }

  .input-section .textarea {
    min-height: 120px;
  }

  .loading {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 2rem;
    text-align: center;
    border: 2px dashed #646cff;
    border-radius: 8px;
    background-color: #f8f9ff;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #e0e0e0;
    border-top: 4px solid #646cff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    flex-shrink: 0;
  }

  .loading-text {
    flex: 1;
    text-align: left;
  }

  .loading-text p {
    margin: 0.25rem 0;
  }

  .loading-text small {
    color: #888;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>