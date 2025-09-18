<script>
  export let analysisResult = null;

  let svgWidth = 800;
  let svgHeight = 400;
  let margin = { top: 60, right: 40, bottom: 80, left: 40 };
  let chartWidth = svgWidth - margin.left - margin.right;
  let chartHeight = svgHeight - margin.top - margin.bottom;

  $: timelineData = analysisResult ? processTimelineData(analysisResult) : null;

  // Convert IP address to integer for calculations
  function ipToInt(ip) {
    const parts = ip.split('.').map(Number);
    return (parts[0] << 24) + (parts[1] << 16) + (parts[2] << 8) + parts[3];
  }

  // Convert integer back to IP address
  function intToIp(int) {
    return [
      (int >>> 24) & 255,
      (int >>> 16) & 255,
      (int >>> 8) & 255,
      int & 255
    ].join('.');
  }

  // Process analysis data into timeline format
  function processTimelineData(result) {
    if (!result.valid_cidrs || result.valid_cidrs.length === 0) {
      return null;
    }

    // Create ranges for valid CIDRs
    let ranges = result.valid_cidrs.map(cidr => {
      const startInt = ipToInt(cidr.network);
      const endInt = ipToInt(cidr.broadcast);
      
      // Determine type based on category
      let type = 'allocated';
      if (cidr.category === 'vpc') {
        type = 'vpc';
      } else if (cidr.category === 'subnet') {
        type = 'subnet';
      }
      
      return {
        type: type,
        start: startInt,
        end: endInt,
        label: cidr.original,
        category: cidr.category || 'other',
        size: cidr.total_ips,
        startIp: cidr.network,
        endIp: cidr.broadcast
      };
    });

    // Add gaps
    if (result.gaps) {
      result.gaps.forEach(gap => {
        ranges.push({
          type: 'gap',
          start: ipToInt(gap.start_ip),
          end: ipToInt(gap.end_ip),
          label: `Gap: ${gap.size.toLocaleString()} IPs`,
          size: gap.size,
          startIp: gap.start_ip,
          endIp: gap.end_ip,
          suggestedCidr: gap.suggested_cidr
        });
      });
    }

    // Mark overlapped ranges
    if (result.overlaps && result.overlaps.length > 0) {
      const overlappedCidrs = new Set();
      result.overlaps.forEach(overlap => {
        overlappedCidrs.add(overlap.cidr1);
        overlappedCidrs.add(overlap.cidr2);
      });
      
      ranges.forEach(range => {
        if (range.type === 'allocated' && overlappedCidrs.has(range.label)) {
          range.type = 'overlap';
        }
      });
    }

    // Sort ranges by start IP
    ranges.sort((a, b) => a.start - b.start);

    // Find the overall range
    const minIp = Math.min(...ranges.map(r => r.start));
    const maxIp = Math.max(...ranges.map(r => r.end));
    const totalRange = maxIp - minIp;

    // Calculate positions and widths for visualization
    ranges.forEach(range => {
      range.x = ((range.start - minIp) / totalRange) * chartWidth;
      range.width = ((range.end - range.start + 1) / totalRange) * chartWidth;
      range.width = Math.max(range.width, 2); // Minimum width for visibility
    });

    return {
      ranges,
      minIp: intToIp(minIp),
      maxIp: intToIp(maxIp),
      totalRange,
      totalIps: totalRange + 1
    };
  }

  // Format large numbers
  function formatNumber(num) {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M';
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
  }

  // Get color for range type
  function getRangeColor(type) {
    switch (type) {
      case 'allocated': return '#10b981'; // green
      case 'vpc': return '#3b82f6'; // blue
      case 'subnet': return '#10b981'; // green
      case 'gap': return '#f59e0b'; // amber
      case 'overlap': return '#ef4444'; // red
      default: return '#6b7280'; // gray
    }
  }

  // Handle hover events
  let hoveredRange = null;
  let mouseX = 0;
  let mouseY = 0;

  function handleMouseMove(event, range) {
    hoveredRange = range;
    mouseX = event.clientX;
    mouseY = event.clientY;
  }

  function handleMouseLeave() {
    hoveredRange = null;
  }
</script>

<div class="timeline-container">
  <h3>IP Address Space Timeline</h3>
  
  {#if timelineData}
    <div class="timeline-info">
      <div class="info-item">
        <span class="label">Range:</span>
        <span class="value">{timelineData.minIp} - {timelineData.maxIp}</span>
      </div>
      <div class="info-item">
        <span class="label">Total IPs:</span>
        <span class="value">{formatNumber(timelineData.totalIps)}</span>
      </div>
    </div>

    <div class="timeline-svg-container">
      <svg width={svgWidth} height={svgHeight}>
        <!-- Background -->
        <rect 
          x={margin.left} 
          y={margin.top} 
          width={chartWidth} 
          height={chartHeight} 
          fill="#1a1a1a" 
          stroke="#333" 
          stroke-width="1"
        />

        <!-- Timeline ranges -->
        {#each timelineData.ranges as range, i}
          <g>
            <!-- Range rectangle -->
            <rect
              x={margin.left + range.x}
              y={margin.top + (range.category === 'vpc' ? 10 : range.category === 'subnet' ? 30 : 20)}
              width={range.width}
              height={range.category === 'vpc' ? chartHeight - 20 : range.category === 'subnet' ? chartHeight - 60 : chartHeight - 40}
              fill={getRangeColor(range.type)}
              stroke={range.category === 'vpc' ? '#3b82f6' : '#333'}
              stroke-width={range.category === 'vpc' ? '2' : '1'}
              opacity={range.category === 'vpc' ? '0.3' : '0.8'}
              on:mousemove={(e) => handleMouseMove(e, range)}
              on:mouseleave={handleMouseLeave}
              class="range-rect"
              role="button"
              tabindex="0"
              aria-label={`${range.label}: ${range.startIp} to ${range.endIp}, ${range.size.toLocaleString()} IPs`}
            />

            <!-- Range label (if wide enough) -->
            {#if range.width > 80}
              <text
                x={margin.left + range.x + range.width / 2}
                y={margin.top + chartHeight / 2}
                text-anchor="middle"
                fill="white"
                font-size="12"
                font-weight="500"
              >
                {range.label}
              </text>
              <text
                x={margin.left + range.x + range.width / 2}
                y={margin.top + chartHeight / 2 + 16}
                text-anchor="middle"
                fill="white"
                font-size="10"
                opacity="0.8"
              >
                {formatNumber(range.size)} IPs
              </text>
            {/if}
          </g>
        {/each}

        <!-- Timeline axis -->
        <line
          x1={margin.left}
          y1={margin.top + chartHeight + 10}
          x2={margin.left + chartWidth}
          y2={margin.top + chartHeight + 10}
          stroke="#666"
          stroke-width="2"
        />

        <!-- Start IP label -->
        <text
          x={margin.left}
          y={margin.top + chartHeight + 30}
          text-anchor="start"
          fill="#888"
          font-size="12"
        >
          {timelineData.minIp}
        </text>

        <!-- End IP label -->
        <text
          x={margin.left + chartWidth}
          y={margin.top + chartHeight + 30}
          text-anchor="end"
          fill="#888"
          font-size="12"
        >
          {timelineData.maxIp}
        </text>

        <!-- Title -->
        <text
          x={svgWidth / 2}
          y={30}
          text-anchor="middle"
          fill="#646cff"
          font-size="16"
          font-weight="600"
        >
          IP Address Space Visualization
        </text>
      </svg>
    </div>

    <!-- Legend -->
    <div class="legend">
      <div class="legend-item">
        <div class="legend-color vpc"></div>
        <span>VPC CIDR Blocks</span>
      </div>
      <div class="legend-item">
        <div class="legend-color subnet"></div>
        <span>Subnet CIDRs</span>
      </div>
      <div class="legend-item">
        <div class="legend-color gap"></div>
        <span>Available IP Gaps</span>
      </div>
      <div class="legend-item">
        <div class="legend-color overlap"></div>
        <span>Overlapping Ranges</span>
      </div>
    </div>

    <!-- Tooltip -->
    {#if hoveredRange}
      <div 
        class="tooltip" 
        style="left: {mouseX + 10}px; top: {mouseY - 60}px;"
      >
        <div class="tooltip-title">{hoveredRange.label}</div>
        <div class="tooltip-detail">Start: {hoveredRange.startIp}</div>
        <div class="tooltip-detail">End: {hoveredRange.endIp}</div>
        <div class="tooltip-detail">Size: {hoveredRange.size.toLocaleString()} IPs</div>
        {#if hoveredRange.suggestedCidr}
          <div class="tooltip-detail">Suggested: {hoveredRange.suggestedCidr}</div>
        {/if}
      </div>
    {/if}

  {:else}
    <div class="empty-state">
      <p>No CIDR data to visualize. Analyze some CIDR ranges to see the timeline.</p>
    </div>
  {/if}
</div>

<style>
  .timeline-container {
    background: #1a1a1a;
    border-radius: 8px;
    border: 1px solid #333;
    padding: 1.5rem;
    margin-top: 1rem;
  }

  .timeline-container h3 {
    margin: 0 0 1rem 0;
    color: #646cff;
  }

  .timeline-info {
    display: flex;
    gap: 2rem;
    margin-bottom: 1rem;
    font-size: 0.875rem;
  }

  .info-item {
    display: flex;
    gap: 0.5rem;
  }

  .label {
    color: #888;
  }

  .value {
    color: white;
    font-weight: 500;
  }

  .timeline-svg-container {
    overflow-x: auto;
    border-radius: 4px;
  }

  .range-rect {
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .range-rect:hover {
    opacity: 1 !important;
  }

  .legend {
    display: flex;
    gap: 1.5rem;
    margin-top: 1rem;
    font-size: 0.875rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .legend-color {
    width: 16px;
    height: 16px;
    border-radius: 3px;
  }

  .legend-color.vpc {
    background-color: #3b82f6;
  }

  .legend-color.subnet {
    background-color: #10b981;
  }

  .legend-color.gap {
    background-color: #f59e0b;
  }

  .legend-color.overlap {
    background-color: #ef4444;
  }

  .tooltip {
    position: fixed;
    background: #2a2a2a;
    border: 1px solid #444;
    border-radius: 6px;
    padding: 0.75rem;
    font-size: 0.875rem;
    z-index: 1000;
    pointer-events: none;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  .tooltip-title {
    font-weight: 600;
    color: #646cff;
    margin-bottom: 0.5rem;
  }

  .tooltip-detail {
    color: #ccc;
    margin-bottom: 0.25rem;
  }

  .tooltip-detail:last-child {
    margin-bottom: 0;
  }

  .empty-state {
    text-align: center;
    padding: 2rem;
    color: #888;
  }

  .empty-state p {
    margin: 0;
  }
</style>