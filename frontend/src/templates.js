export function initTemplates() {
    function field(label, name, type, placeholder, opts) {
        opts = opts || {};
        let html = '<label>' + label + '</label>';
        if (type === 'select') {
            html += '<select name="' + name + '">';
            (opts.options || []).forEach(function (o) {
                html += '<option value="' + o.value + '">' + o.text + '</option>';
            });
            html += '</select>';
        } else if (type === 'checkbox') {
            html += '<input type="checkbox" name="' + name + '"';
            if (opts.checked) html += ' checked';
            html += '>';
        } else {
            html += '<input type="' + type + '" name="' + name + '" placeholder="' + (placeholder || '') + '"';
            if (opts.required) html += ' required';
            html += '>';
        }
        return html;
    }

    return {
        text: function () {
            return '<div class="qr-form">' +
                field('Text', 'text', 'text', 'Enter text to encode') +
                '</div>';
        },

        url: function () {
            return '<div class="qr-form">' +
                field('URL', 'url', 'url', 'https://example.com') +
                '</div>';
        },

        email: function () {
            return '<div class="qr-form">' +
                field('To', 'to', 'email', 'addr@example.com') +
                field('Subject', 'subject', 'text', 'Subject') +
                field('Body', 'body', 'text', 'Email body') +
                '</div>';
        },

        phone: function () {
            return '<div class="qr-form">' +
                field('Phone', 'phone', 'tel', '+1234567890') +
                '</div>';
        },

        sms: function () {
            return '<div class="qr-form">' +
                field('Phone', 'phone', 'tel', '+1234567890') +
                field('Message', 'message', 'text', 'Message text') +
                '</div>';
        },

        wifi: function () {
            return '<div class="qr-form">' +
                field('SSID', 'ssid', 'text', 'Network name') +
                field('Password', 'password', 'text', 'WiFi password') +
                field('Encryption', 'encryption', 'select', '', {
                    options: [
                        {value: 'WPA', text: 'WPA/WPA2'},
                        {value: 'WEP', text: 'WEP'},
                        {value: 'nopass', text: 'None'}
                    ]
                }) +
                field('Hidden', 'hidden', 'checkbox', '', {checked: false}) +
                '</div>';
        },

        vcard: function () {
            return '<div class="qr-form">' +
                field('Name', 'name', 'text', 'Full name') +
                field('Phone', 'phone', 'tel', '+1234567890') +
                field('Email', 'email', 'email', 'addr@example.com') +
                field('Organization', 'org', 'text', 'Company') +
                field('Title', 'title', 'text', 'Job title') +
                field('Address', 'address', 'text', 'Street, City') +
                '</div>';
        },

        geo: function () {
            return '<div class="qr-form">' +
                field('Latitude', 'latitude', 'number', '52.5200') +
                field('Longitude', 'longitude', 'number', '13.4050') +
                '</div>';
        },

        calendar: function () {
            return '<div class="qr-form">' +
                field('Title', 'title', 'text', 'Event title') +
                field('Start', 'start', 'text', 'YYYYMMDDTHHMMSS') +
                field('End', 'end', 'text', 'YYYYMMDDTHHMMSS') +
                field('Location', 'location', 'text', 'Event location') +
                field('Description', 'description', 'text', 'Event description') +
                '</div>';
        },

        bitcoin: function () {
            return '<div class="qr-form">' +
                field('Address', 'address', 'text', 'Bitcoin address') +
                field('Amount', 'amount', 'number', '0.01') +
                field('Label', 'label', 'text', 'Payment label') +
                '</div>';
        }
    };
}
