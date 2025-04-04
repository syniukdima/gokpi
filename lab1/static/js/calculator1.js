"use strict";

var calculator1Form = document.getElementById('calculator1Form');
calculator1Form.addEventListener('submit', function(e) {
    e.preventDefault();

    var formData = {
        h: parseFloat(document.getElementById('inputH').value),
        c: parseFloat(document.getElementById('inputC').value),
        s: parseFloat(document.getElementById('inputS').value),
        n: parseFloat(document.getElementById('inputN').value),
        o: parseFloat(document.getElementById('inputO').value),
        w: parseFloat(document.getElementById('inputW').value),
        a: parseFloat(document.getElementById('inputA').value)
    };

    fetch('/api/calculate1', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
    })
        .then(function(response) {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(function(result) {
            var dryMassHtml =
                '<p>Водень (H): ' + result.dryMass.h.toFixed(2) + '%</p>' +
                '<p>Вуглець (C): ' + result.dryMass.c.toFixed(2) + '%</p>' +
                '<p>Сірка (S): ' + result.dryMass.s.toFixed(2) + '%</p>' +
                '<p>Азот (N): ' + result.dryMass.n.toFixed(2) + '%</p>' +
                '<p>Кисень (O): ' + result.dryMass.o.toFixed(2) + '%</p>' +
                '<p>Зола (A): ' + result.dryMass.a.toFixed(2) + '%</p>';

            document.getElementById('dryMassResults').innerHTML = dryMassHtml;

            var combustibleMassHtml =
                '<p>Водень (H): ' + result.combustibleMass.h.toFixed(2) + '%</p>' +
                '<p>Вуглець (C): ' + result.combustibleMass.c.toFixed(2) + '%</p>' +
                '<p>Сірка (S): ' + result.combustibleMass.s.toFixed(2) + '%</p>' +
                '<p>Азот (N): ' + result.combustibleMass.n.toFixed(2) + '%</p>' +
                '<p>Кисень (O): ' + result.combustibleMass.o.toFixed(2) + '%</p>';

            document.getElementById('combustibleMassResults').innerHTML = combustibleMassHtml;

            var heatingValueHtml =
                '<p>Нижча робоча теплота згоряння: ' + result.heatingValue.working.toFixed(1) + ' кДж/кг</p>' +
                '<p>Нижча суха теплота згоряння: ' + result.heatingValue.dry.toFixed(1) + ' кДж/кг</p>' +
                '<p>Нижча горюча теплота згоряння: ' + result.heatingValue.combustible.toFixed(1) + ' кДж/кг</p>';

            document.getElementById('heatingValueResults').innerHTML = heatingValueHtml;
        })
        .catch(function(error) {
            window.alert('Помилка при розрахунках. Перевірте введені дані.');
            if (window.console && window.console.error) {
                window.console.error('Error:', error);
            }
        });
});