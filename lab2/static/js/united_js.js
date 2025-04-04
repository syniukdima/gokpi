"use strict";

/* jshint esversion: 6 */
/* eslint-env es6 */

function calculate() {
    const errorDiv = document.getElementById('error');
    errorDiv.style.display = 'none';

    const coalVolume = parseFloat(document.getElementById('coalVolume').value) || 0;
    const oilFuelVolume = parseFloat(document.getElementById('oilFuelVolume').value) || 0;
    const naturalGasVolume = parseFloat(document.getElementById('naturalGasVolume').value) || 0;

    fetch('http://localhost:8080/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            coalVolume: coalVolume,
            oilFuelVolume: oilFuelVolume,
            naturalGasVolume: naturalGasVolume
        })
    })
        .then(function(response) {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(function(data) {
            document.getElementById('coalSolidParticlesEmission').textContent =
                'Емісія твердих частинок при спалюванні: ' + data.coal.solidParticlesEmission.toFixed(2);
            document.getElementById('coalGrossEmission').textContent =
                'Валовий викид при спалюванні: ' + data.coal.grossEmission.toFixed(2);

            document.getElementById('oilFuelSolidParticlesEmission').textContent =
                'Емісія твердих частинок при спалюванні: ' + data.oilFuel.solidParticlesEmission.toFixed(2);
            document.getElementById('oilFuelGrossEmission').textContent =
                'Валовий викид при спалюванні: ' + data.oilFuel.grossEmission.toFixed(2);

            document.getElementById('naturalGasSolidParticlesEmission').textContent =
                'Емісія твердих частинок при спалюванні: ' + data.naturalGas.solidParticlesEmission.toFixed(2);
            document.getElementById('naturalGasGrossEmission').textContent =
                'Валовий викид при спалюванні: ' + data.naturalGas.grossEmission.toFixed(2);
        })
        .catch(function(error) {
            console.error('Error:', error);
            errorDiv.style.display = 'block';
        });
}