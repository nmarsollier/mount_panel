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
        motorsEnabled: true,
        motorsSwitch: false,
        debug: {},
        serverSettings: {lat: 0, lon: 0, elevation: 0},
        settingsDraft: {lat: 0, lon: 0, elevation: 0},
        settingsDirty: false,
        mountTime: '--',
        slew: {degrees: 5, speed: 4},
        joySpeed: -1,        // index into moveSpeeds; set to fastest once fetched
        moveSpeeds: [],      // deg/s ladder, from capabilities
        slewTarget: {ra: {h: 0, m: 0, s: 0}, dec: {d: 0, m: 0, s: 0}, speed: 4},
        serverTracking: 'none',
        isParked: false,
        selectedTracking: null,
        trackingOptions: ['none', 'sidereal', 'lunar', 'solar'],
        joyRates: {ra: 0, dec: 0},
        eaf: {detected: false, power: false, health: 'unknown',
              focuser: {counter_steps: 0, counter_valid: true, is_moving: false, enabled: true},
              fan: {mode: 'off', power: false, speed: 0, rpm: null},
              mirror: {error: null, mirror_c: null},
              environment: {error: null, ambient_c: null, humidity_pct: null, dew_point_c: null},
              vibration: {error: null, rms_total: null, accel_peak_total: null, vel_rms_total: null, vel_peak_total: null, iso_zone: null},
              orientation: {error: null, pitch_deg: null, roll_deg: null, ax_g: null, ay_g: null, az_g: null}},
        focusStep: 100,
        vibHistory: { samples: [], stats: {} },
        moveHistory: { intervals: [] },
        mountMoving: false,
        vibHover: null,

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
                this.motorsEnabled = (j.motors_enabled !== false);
                this.motorsSwitch = (j.motors_switch === true);
                // "slewing" only: the mount reports a nonzero speed while
                // sidereal tracking, which is not a movement event.
                this.mountMoving = (j.status === 'slewing');
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
        toggleAmb() {
            if (this.eaf.amb_detected) {
                this.run(app().DisconnectAmb());
            } else {
                this.run(app().ConnectAmb());
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

        toggleMotors() {
            this.run(app().SetMotorsEnabled(!this.motorsEnabled));
        },

        // --- Joystick (continuous move) ---
        joyStart(axis, dir) {
            const dps = this.moveSpeeds[this.joySpeed] || 0;
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

        /* --- Focuser (EAF: focuser + fan) --- */
        fetchEAFStatus() {
            app().EAFStatus().then(j => {
                if (j) this.eaf = j;
                this.normalizeTab();
            }).catch(() => {
            });
        },

        /* --- Vibration history --- */
        fetchVibrationHistory() {
            app().VibrationHistory().then(j => {
                if (j) this.vibHistory = j;
            }).catch(() => {
            });
        },

        /* --- Mount movement history (chart overlay) --- */
        fetchMovementHistory() {
            app().MovementHistory().then(j => {
                if (j) this.moveHistory = j;
            }).catch(() => {
            });
        },

        /* Load the MoveAxis speed ladder from the mount capabilities. */
        fetchMoveSpeeds() {
            app().MoveAxisSpeeds().then(speeds => {
                if (speeds && speeds.length) {
                    this.moveSpeeds = speeds;
                    if (this.joySpeed < 0 || this.joySpeed >= speeds.length) {
                        this.joySpeed = speeds.length - 1;
                    }
                }
            }).catch(() => {
            });
        },

        /* Move to the first visible tab when the active one's device disconnects. */
        normalizeTab() {
            const visible = [];
            if (this.mountConnected) visible.push('principal', 'movimiento', 'settings');
            if (this.eaf.detected) visible.push('enfoque');
            if (this.eaf.fan_detected) visible.push('fan');
            if (this.eaf.amb_detected) visible.push('ambiente', 'vibraciones');
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
            if (this.eaf.health === 'fault') return 'Fault';
            if (!this.eaf.focuser.enabled) return 'Disabled';
            if (this.eaf.focuser.is_moving) return 'Moving';
            return 'Ok';
        },

        /* The firmware reports counter_valid=false after a disable, power loss
         * or reboot — N.I.N.A. must not autofocus until the focuser is
         * re-enabled. */
        focusReferenceMissing() {
            return this.eaf.detected && this.eaf.focuser.counter_valid === false;
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
            if (!this.eaf.detected) return 'Not detected — tap "EAF" above to connect';
            if (this.eaf.health === 'fault') return 'Fault — check the focuser driver';
            if (!this.eaf.focuser.enabled) return 'Disabled';
            if (this.eaf.focuser.is_moving) return 'Moving';
            return 'Ok';
        },

        eafStatusClass() {
            if (!this.eaf.detected || this.eaf.health === 'fault') return 'status-error';
            if (!this.eaf.focuser.enabled) return 'status-warn';
            return 'status-ok';
        },

        sensorOk(key) {
            if (key === 'ambiente') return this.eaf.amb_detected && this.eaf.environment.ambient_c != null;
            if (key === 'espejo') return this.eaf.fan_detected && this.eaf.mirror.mirror_c != null;
            if (key === 'vibracion') return this.eaf.amb_detected && this.eaf.vibration.error == null;
            return false;
        },

        sensorError() {
            if (!this.eaf.amb_detected) {
                return 'Not detected — tap "AMB" above to connect';
            }
            const errs = [];
            if (this.eaf.environment.error) errs.push(this.eaf.environment.error);
            if (this.eaf.vibration.error) errs.push(this.eaf.vibration.error);
            return errs.join(' · ');
        },

        /* --- Vibration chart (inline SVG, no library) --- */
        vibSamples() {
            return this.vibHistory.samples || [];
        },
        // Fixed geometry: viewBox 0 0 600 220; y maps 0..8 mm/s (ISO 10816
        // Class II boundaries at 1.12 / 2.8 / 7.1), x maps the last 5 min on
        // the shared time axis (same as the movement strip below).
        vibY(v) {
            const top = 12, bottom = 20, H = 220, yMax = 8;
            return top + (1 - Math.min(v, yMax) / yMax) * (H - top - bottom);
        },
        chartX(ts) {
            const W = 600, pad = 4, windowMs = 300000;
            const now = Date.now();
            return pad + (ts - (now - windowMs)) / windowMs * (W - 2 * pad);
        },
        vibPoints() {
            const samples = this.vibSamples();
            if (samples.length < 2) return '';
            return samples.map(s => `${this.chartX(s.ts_ms).toFixed(1)},${this.vibY(s.vel_rms).toFixed(1)}`).join(' ');
        },
        vibHoverAt(e) {
            const samples = this.vibSamples();
            if (samples.length < 2) return;
            const rect = e.currentTarget.getBoundingClientRect();
            const frac = (e.clientX - rect.left) / rect.width;
            const idx = Math.max(0, Math.min(samples.length - 1, Math.round(frac * (samples.length - 1))));
            this.vibHover = samples[idx];
        },
        vibHoverClear() {
            this.vibHover = null;
        },
        vibHoverText() {
            if (!this.vibHover) return '';
            const d = new Date(this.vibHover.ts_ms);
            const pad = n => String(n).padStart(2, '0');
            return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())} · ${this.vibHover.vel_rms.toFixed(3)} mm/s`;
        },
        moveBands() {
            const intervals = this.moveHistory.intervals || [];
            if (!intervals.length) return [];
            const W = 600, pad = 4;
            return intervals.map(iv => {
                let x1 = this.chartX(iv.start_ms);
                let x2 = this.chartX(iv.end_ms);
                x1 = Math.max(pad, Math.min(W - pad, x1));
                x2 = Math.max(pad, Math.min(W - pad, x2));
                // A very brief move (one sample) collapses to zero width; give
                // it a ~2px sliver so it is still visible.
                if (x2 - x1 < 2) x2 = Math.min(W - pad, x1 + 2);
                return { left: (x1 / W) * 100, width: ((x2 - x1) / W) * 100, start: iv.start_ms };
            }).filter(b => b.width > 0);
        },

        init() {
            this.fetchStatus();
            this.fetchEAFStatus();
            this.fetchMoveSpeeds();
            this.fetchVibrationHistory();
            this.fetchMovementHistory();
            setInterval(() => { this.fetchStatus(); this.fetchEAFStatus(); this.fetchMoveSpeeds(); this.fetchVibrationHistory(); this.fetchMovementHistory(); }, 1000);
        }
    }
}
