function mountApp() {
    // Wails exposes the Go App binding at window.go.<package>.<Struct>.
    const app = () => window.go.ui.App;

    return {
        tab: 'principal',
        status: {ra: '--:--:--', dec: '--:--:--'},
        lst: '--:--:--',
        pierSide: '--',
        statusText: 'UNKNOWN',
        isError: false,
        isDisconnected: false,
        hasPower: true,
        debug: {},
        serverSettings: {lat: 0, lon: 0, elevation: 0},
        settingsDraft: {lat: 0, lon: 0, elevation: 0},
        settingsDirty: false,
        mountTime: '--',
        timeAutoSet: false,
        errorMessage: '',
        errorTimer: null,
        slew: {degrees: 5, speed: 4},
        joySpeed: 4,
        slewTarget: {ra: {h: 0, m: 0, s: 0}, dec: {d: 0, m: 0, s: 0}, speed: 4},
        serverTracking: 'none',
        isHome: false,
        isParked: false,
        selectedTracking: null,
        trackingOptions: ['none', 'sidereal', 'lunar', 'solar'],
        joyRates: {ra: 0, dec: 0},

        /* Show an error message for 5 seconds, then auto-clear. */
        showError(msg) {
            this.errorMessage = msg;
            if (this.errorTimer) clearTimeout(this.errorTimer);
            this.errorTimer = setTimeout(() => {
                this.errorMessage = '';
            }, 5000);
        },

        /* Run a binding promise, refresh status, and surface errors. */
        run(promise) {
            return promise
                .then(() => this.fetchStatus())
                .catch(e => this.showError(String(e)));
        },

        fetchStatus() {
            app().Status().then(j => {
                this.statusText = j.status || 'UNKNOWN';
                this.status = {ra: j.ra || '--:--:--', dec: j.dec || '--:--:--'};
                this.lst = j.lst || '--:--:--';
                this.pierSide = j.pier_side || '--';
                this.mountTime = j.time || '--';
                this.serverTracking = j.tracking || 'none';
                this.isHome = j.is_home || false;
                this.isParked = (j.status === 'parked');
                this.isError = (j.status === 'error');
                this.isDisconnected = (j.status === 'disconnected');
                this.hasPower = (j.power !== false);
                this.debug = j.debug || {};

                const s = j.settings;
                if (s) {
                    this.serverSettings = {
                        lat: s.lat != null ? s.lat : 0,
                        lon: s.lon != null ? s.lon : 0,
                        elevation: s.elevation != null ? s.elevation : 0,
                    };
                    // Only overwrite the draft when the user is not editing.
                    if (!this.settingsDirty) {
                        this.settingsDraft = {
                            lat: this.serverSettings.lat,
                            lon: this.serverSettings.lon,
                            elevation: this.serverSettings.elevation,
                        };
                    }
                }

                if (this.selectedTracking === null) {
                    const idx = this.trackingOptions.indexOf(j.tracking);
                    this.selectedTracking = idx >= 0 ? idx : 0;
                }
            }).catch(() => {
            });
        },

        stop() {
            this.run(app().Stop());
        },
        home() {
            this.run(app().Home());
        },
        park() {
            this.run(app().Park());
        },
        unpark() {
            this.run(app().Unpark());
        },

        reset() {
            if (confirm('Reboot the mount controller?')) {
                this.run(app().Reset());
            }
        },

        setTracking() {
            this.run(app().SetTracking(this.trackingOptions[this.selectedTracking]));
        },

        setLimit(action) {
            this.run(app().SetLimit(action));
        },

        // --- Joystick (continuous move) ---
        joyStart(axis, dir) {
            const speeds = {1: 1.0, 2: 3.0, 3: 4.5, 4: 6.0};
            const dps = speeds[this.joySpeed] || 6.0;
            this.joyRates[axis] = dir * dps;
            app().MoveAxisSpeed(this.joyRates.ra, this.joyRates.dec)
                .catch(e => this.showError(String(e)));
        },

        joyStop(axis) {
            this.joyRates[axis] = 0;
            app().MoveAxisSpeed(this.joyRates.ra, this.joyRates.dec)
                .catch(e => this.showError(String(e)));
        },

        doSlew(axis, dir) {
            const axisName = axis === 'ra' ? 'ra' : 'dec';
            const deg = Math.min(180, Math.max(1, Math.round(Math.abs(this.slew.degrees || 5))));
            const spd = this.slew.speed || 4;
            this.run(app().MoveAxisBy(axisName, dir * deg, spd));
        },

        doSlewToCoordinates() {
            const rh = this.slewTarget.ra.h;
            const rm = this.slewTarget.ra.m;
            const rs = this.slewTarget.ra.s;
            const dd = this.slewTarget.dec.d;
            const dm = this.slewTarget.dec.m;
            const ds = this.slewTarget.dec.s;

            if (isNaN(rh) || rh < 0 || rh > 23 || isNaN(rm) || rm < 0 || rm > 59 || isNaN(rs) || rs < 0 || rs >= 60) {
                this.showError('Invalid RA — H: 0–23, M: 0–59, S: 0–59.999');
                return;
            }
            if (isNaN(dd) || dd < -180 || dd > 180 || isNaN(dm) || dm < 0 || dm > 59 || isNaN(ds) || ds < 0 || ds >= 60) {
                this.showError('Invalid DEC — D: −180…+180, M: 0–59, S: 0–59.999');
                return;
            }

            const raHours = rh + rm / 60 + rs / 3600;
            const absDec = Math.abs(dd) + dm / 60 + ds / 3600;
            const decDeg = dd < 0 ? -absDec : absDec;

            this.run(app().SlewTo(raHours, decDeg, this.slewTarget.speed));
        },

        /* Read and validate the location draft; null when invalid. */
        locationDraft() {
            const d = this.settingsDraft;
            const lat = parseFloat(d.lat), lon = parseFloat(d.lon), elev = parseFloat(d.elev);
            if (isNaN(lat) || isNaN(lon) || isNaN(elev)) {
                this.showError('Location fields must be valid numbers');
                return null;
            }
            return {lat, lon, elev};
        },

        updateTime() {
            const loc = this.locationDraft();
            if (!loc) return;
            this.run(
                app().SetTime(new Date().toISOString())
                    .then(() => app().SetLocation(loc.lat, loc.lon, loc.elev))
                    .then(() => { this.settingsDirty = false; })
            );
        },

        updateSettings() {
            const loc = this.locationDraft();
            if (!loc) return;
            this.run(
                app().SetLocation(loc.lat, loc.lon, loc.elev)
                    .then(() => { this.settingsDirty = false; })
            );
        },

        fmtUptime(s) {
            if (s == null) return '—';
            const h = Math.floor(s / 3600);
            const m = Math.floor((s % 3600) / 60);
            const sec = s % 60;
            if (h > 0) return h + 'h ' + m + 'm ' + sec + 's';
            if (m > 0) return m + 'm ' + sec + 's';
            return sec + 's';
        },

        fatalMessage() {
            if (!this.hasPower) {
                return 'Motors unpowered — turn on the 12V switch. Movement is disabled.';
            }
            return 'Mount error — reboot required. Mount will not accept movement commands.';
        },

        disconnectedMessage() {
            return 'Mount not connected — plug in and power the mount. Reconnecting…';
        },

        init() {
            this.fetchStatus();
            setInterval(() => this.fetchStatus(), 1000);
        }
    }
}
