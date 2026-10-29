function mountApp() {
    // Wails exposes the Go App binding at window.go.<package>.<Struct>.
    const app = () => window.go.ui.App;

    return {
        tab: 'principal',
        status: {ra: '--:--:--', dec: '--:--:--'},
        lst: '--:--:--',
        pierSide: '--',
        isError: false,
        isDisconnected: false,
        mountConnected: false,
        hasPower: true,
        debug: {},
        serverSettings: {lat: 0, lon: 0, elevation: 0},
        settingsDraft: {lat: 0, lon: 0, elevation: 0},
        settingsDirty: false,
        mountTime: '--',
        slew: {degrees: 5, speed: 4},
        joySpeed: 4,
        slewTarget: {ra: {h: 0, m: 0, s: 0}, dec: {d: 0, m: 0, s: 0}, speed: 4},
        serverTracking: 'none',
        isParked: false,
        selectedTracking: null,
        trackingOptions: ['none', 'sidereal', 'lunar', 'solar'],
        joyRates: {ra: 0, dec: 0},
        eaf: {detected: false, power: false, health: 'unknown',
              focuser: {counter_steps: 0, is_moving: false, enabled: true},
              fan: {mode: 'off', power: false, speed: 0, rpm: null},
              environment: {error: null, ambient_c: null, humidity_pct: null, mirror_c: null, dew_point_c: null, mirror_ambient_delta_c: null},
              vibration: {error: null, rms_total: 0}},
        focusStep: 100,

        /* Run a binding promise and refresh status. */
        run(promise) {
            return promise
                .then(() => this.fetchStatus())
                .then(() => this.fetchEAFStatus())
                .catch(() => {});
        },

        fetchStatus() {
            app().Status().then(j => {
                this.status = {ra: j.ra || '--:--:--', dec: j.dec || '--:--:--'};
                this.lst = j.lst || '--:--:--';
                this.pierSide = j.pier_side || '--';
                this.mountTime = j.time || '--';
                this.serverTracking = j.tracking || 'none';
                this.isParked = (j.status === 'parked');
                this.isError = (j.status === 'error');
                this.isDisconnected = (j.status === 'disconnected');
                this.mountConnected = (j.status !== 'disconnected');
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
                this.normalizeTab();
            }).catch(() => {
            });
        },

        stop() {
            this.run(app().Stop());
        },
        home() {
            this.run(app().Home());
        },

        /* Header chips toggle the connection: connect when down, disconnect when up. */
        toggleMount() {
            if (this.mountConnected) {
                this.run(app().DisconnectMount());
            } else {
                this.run(app().ConnectMount());
            }
        },
        toggleEAF() {
            if (this.eaf.detected) {
                this.run(app().DisconnectEAF());
            } else {
                this.run(app().ConnectEAF());
            }
        },
        toggleFan() {
            if (this.eaf.fan_detected) {
                this.run(app().DisconnectFan());
            } else {
                this.run(app().ConnectFan());
            }
        },
        toggleTemp() {
            if (this.eaf.temp_detected) {
                this.run(app().DisconnectTemp());
            } else {
                this.run(app().ConnectTemp());
            }
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
                .catch(() => {});
        },

        joyStop(axis) {
            this.joyRates[axis] = 0;
            app().MoveAxisSpeed(this.joyRates.ra, this.joyRates.dec)
                .catch(() => {});
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
                console.warn('Invalid RA — H: 0–23, M: 0–59, S: 0–59.999');
                return;
            }
            if (isNaN(dd) || dd < -180 || dd > 180 || isNaN(dm) || dm < 0 || dm > 59 || isNaN(ds) || ds < 0 || ds >= 60) {
                console.warn('Invalid DEC — D: −180…+180, M: 0–59, S: 0–59.999');
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
                console.warn('Location fields must be valid numbers');
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

        /* --- Enfoque (EAF: focuser + fan) --- */
        fetchEAFStatus() {
            app().EAFStatus().then(j => {
                if (j) this.eaf = j;
                this.normalizeTab();
            }).catch(() => {
            });
        },

        /* Move to the first visible tab when the active one's device disconnects. */
        normalizeTab() {
            const visible = [];
            if (this.mountConnected) visible.push('principal', 'movimiento', 'settings');
            if (this.eaf.detected) visible.push('enfoque');
            if (this.eaf.fan_detected) visible.push('fan');
            if (this.eaf.fan_detected || this.eaf.temp_detected) visible.push('sensores');
            if (visible.length && !visible.includes(this.tab)) this.tab = visible[0];
        },

        focusMoveDisabled() {
            if (!this.eaf.detected) return true;
            if (this.eaf.health === 'fault') return true;
            if (!this.eaf.power) return true;
            if (!this.eaf.focuser.enabled) return true;
            if (this.eaf.focuser.is_moving) return true;
            return false;
        },

        focusMove(dir) {
            const step = Math.max(1, Math.round(this.focusStep || 100));
            this.run(app().FocusMove(dir * step));
        },

        focusStop() {
            this.run(app().FocusStop());
        },

        focusToggleEnabled() {
            this.run(app().FocusSetEnabled(!this.eaf.focuser.enabled));
        },

        focuserHeaderState() {
            if (!this.eaf.detected) return '—';
            if (this.eaf.health === 'fault') return 'Fallo';
            if (!this.eaf.focuser.enabled) return 'Deshabilitado';
            if (this.eaf.focuser.is_moving) return 'Moviendo';
            return 'Ok';
        },

        fanSetMode(mode) {
            this.run(app().FanSetMode(mode));
        },

        fanSetSpeed(pct) {
            const v = Math.max(1, Math.min(100, Math.round(pct)));
            this.run(app().FanSetSpeed(v));
        },

        /* Render a sensor reading with fixed decimals, or "—" when absent (null). */
        fmt(v, digits, suffix) {
            if (v == null) return '—';
            return v.toFixed(digits) + suffix;
        },

        fatalMessage() {
            if (!this.hasPower) {
                return 'Motors unpowered — turn on the 12V switch. Movement is disabled.';
            }
            return 'Mount error — reboot required. Mount will not accept movement commands.';
        },

        disconnectedMessage() {
            return 'Mount not connected — plug in and power the mount, then click "Mount" above to connect.';
        },

        mountStatusText() {
            if (this.isError) return this.fatalMessage();
            if (this.isDisconnected) return this.disconnectedMessage();
            return 'Ok';
        },

        mountStatusClass() {
            return (this.isError || this.isDisconnected) ? 'status-error' : 'status-ok';
        },

        eafStatusText() {
            if (!this.eaf.detected) return 'No detectado — tocá "EAF" arriba para conectar';
            if (this.eaf.health === 'fault') return 'Fallo — revisá el driver del enfocador';
            if (!this.eaf.focuser.enabled) return 'Deshabilitado';
            if (this.eaf.focuser.is_moving) return 'Moviendo';
            return 'Ok';
        },

        eafStatusClass() {
            if (!this.eaf.detected || this.eaf.health === 'fault') return 'status-error';
            if (!this.eaf.focuser.enabled) return 'status-warn';
            return 'status-ok';
        },

        sensorOk(key) {
            if (key === 'ambiente') return this.eaf.temp_detected && this.eaf.environment.ambient_c != null;
            if (key === 'espejo') return this.eaf.fan_detected && this.eaf.environment.mirror_c != null;
            if (key === 'vibracion') return this.eaf.temp_detected && this.eaf.vibration.error == null;
            return false;
        },

        sensorError() {
            if (!this.eaf.fan_detected && !this.eaf.temp_detected) {
                return 'No detectado — tocá "Fan" o "Temp" arriba para conectar';
            }
            const errs = [];
            if (this.eaf.environment.error) errs.push(this.eaf.environment.error);
            if (this.eaf.vibration.error) errs.push(this.eaf.vibration.error);
            return errs.join(' · ');
        },

        init() {
            this.fetchStatus();
            this.fetchEAFStatus();
            setInterval(() => { this.fetchStatus(); this.fetchEAFStatus(); }, 1000);
        }
    }
}
