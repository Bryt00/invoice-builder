document.addEventListener('DOMContentLoaded', function() {
    const canvas = document.getElementById('revenueChart');
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    const totalRevStr = canvas.getAttribute('data-revenue');
    const totalRev = parseFloat(totalRevStr) || 0;
    
    const gradient = ctx.createLinearGradient(0, 0, 0, 280);
    gradient.addColorStop(0, 'rgba(245, 158, 11, 0.3)');
    gradient.addColorStop(1, 'rgba(245, 158, 11, 0.0)');

    const data = {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug'],
        datasets: [{
            label: 'Revenue (GH₵)',
            data: [1200, 2400, 3100, 4800, 6200, 8500, 11200, totalRev],
            borderColor: '#f59e0b',
            backgroundColor: gradient,
            borderWidth: 3,
            pointBackgroundColor: '#ffffff',
            pointBorderColor: '#f59e0b',
            pointBorderWidth: 2,
            pointRadius: 4,
            pointHoverRadius: 6,
            fill: true,
            tension: 0.4
        }]
    };

    const config = {
        type: 'line',
        data: data,
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { display: false },
                tooltip: {
                    backgroundColor: '#1f2937',
                    titleFont: { family: 'Nunito Sans', size: 12 },
                    bodyFont: { family: 'Nunito Sans', size: 13, weight: 'bold' },
                    padding: 10,
                    callbacks: {
                        label: function(context) {
                            return 'Revenue: GH₵ ' + context.parsed.y.toFixed(2);
                        }
                    }
                }
            },
            scales: {
                x: {
                    grid: { display: false },
                    ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#6b7280' }
                },
                y: {
                    grid: { color: 'rgba(156, 163, 175, 0.15)', borderDash: [4, 4] },
                    ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#6b7280' },
                    beginAtZero: true
                }
            }
        }
    };

    new Chart(ctx, config);
});
