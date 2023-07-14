#include <vatek_sdk_usbstream.h>

typedef void* hstream_source;
typedef vatek_result(*fpstream_source_start)(hstream_source hsource);
typedef vatek_result(*fpstream_source_check)(hstream_source hsource);
typedef uint8_t*(*fpstream_source_get)(hstream_source hsource);
typedef vatek_result(*fpstream_source_stop)(hstream_source hsource);
typedef void(*fpstream_source_free)(hstream_source hsource);

typedef struct _tsstream_source {
	hstream_source hsource;
	fpstream_source_start start;
	fpstream_source_check check;
	fpstream_source_get get;
	fpstream_source_stop stop;
	fpstream_source_free free;
} tsstream_source, *Ptsstream_source;

typedef struct vatek_context {
	hvatek_devices hdevlist;
	hvatek_chip hchip;
	hvatek_usbstream hustream;
	tsstream_source streamsource;
	hmux_core hmux;
	hmux_channel m_hchannel;
	usbstream_param usbcmd;
} VatekContext;

vatek_result ts_stream_start(hstream_source hsource) {
	return vatek_success;
}

vatek_result ts_stream_stop(hstream_source hsource) {
	return vatek_success;
}

uint8_t* ts_stream_get(hstream_source hsource) {
	return NULL;
}

vatek_result ts_stream_check(hstream_source hsource) {
	return vatek_success;
}

void ts_stream_free(hstream_source hsource) {
}

vatek_result source_sync_get_buffer(void* param, uint8_t** pslicebuf) {
	Ptsstream_source ptssource = (Ptsstream_source)param;
	vatek_result nres = ptssource->check(ptssource->hsource);
	if(nres > vatek_success) {
		*pslicebuf = ptssource->get(ptssource->hsource);
		nres = (vatek_result)1;
	}
	return nres;
}

int GetVatekSDKVersion() {
	return vatek_version();
}

int FreeVatekContext(char* p) {
	vatek_result nres = vatek_success;
	VatekContext* ctx = (VatekContext*)p;
	if(ctx) {
		if(ctx->hustream) {
			nres = vatek_usbstream_stop(ctx->hustream);
			vatek_usbstream_close(ctx->hustream);
		}

		if(ctx->hchip) {
			//reboot chip
			vatek_device_close_reboot(ctx->hchip);
		}
		if(ctx->hdevlist)vatek_device_list_free(ctx->hdevlist);
		if(ctx->streamsource.hsource)
			ctx->streamsource.free(ctx->streamsource.hsource);
		
		free(ctx);
	}
	return nres;
}

char* NewVatekContext() {
	VatekContext* ctx = (VatekContext*)malloc(sizeof(VatekContext));
	memset(ctx, 0, sizeof(VatekContext));
	ctx->usbcmd.mode = ustream_mode_sync;
	ctx->usbcmd.remux = ustream_remux_pcr;
	ctx->usbcmd.pcradjust = pcr_adjust;
	ctx->usbcmd.r2param.freqkhz = 473000;
	ctx->usbcmd.r2param.mode = r2_cntl_path_0;

	ctx->usbcmd.modulator.bandwidth_symbolrate = 6;
	ctx->usbcmd.modulator.type = modulator_dvb_t;
	ctx->usbcmd.modulator.ifmode = ifmode_disable;
	ctx->usbcmd.modulator.mod.dvb_t.constellation = dvb_t_qam64;
	ctx->usbcmd.modulator.mod.dvb_t.fft = fft_8k;
	ctx->usbcmd.modulator.mod.dvb_t.guardinterval = guard_interval_1_16;
	ctx->usbcmd.modulator.mod.dvb_t.coderate = coderate_5_6;

	modulator_param_reset(modulator_atsc, &ctx->usbcmd.modulator);
	ctx->streamsource.hsource = NULL;
	ctx->streamsource.start = ts_stream_start;
	ctx->streamsource.stop = ts_stream_stop;
	ctx->streamsource.get = ts_stream_get;
	ctx->streamsource.check = ts_stream_check;
	ctx->streamsource.free = ts_stream_free;

	return (char*)ctx;
}

int VatekDeviceOpen(char* p) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx) {
		vatek_result nres = vatek_device_list_enum(DEVICE_BUS_USB, service_transform, &ctx->hdevlist);
		if(is_vatek_success(nres)) {
			if(nres == 0) {
				nres = vatek_nodevice;
			} else {
				nres = vatek_device_open(ctx->hdevlist, 0, &ctx->hchip);
			}
		}
		return nres;
	}
	return vatek_memfail;
}

int GetVatekDeviceChipInfo(char* p, int* status, uint32_t* fwVer, int* chipId, uint32_t* service, uint32_t* in, uint32_t* out, uint32_t* peripheral) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx && ctx->hchip) {
		Pchip_info pinfo = vatek_device_get_info(ctx->hchip);
		if(status) *status = pinfo->status;
		if(fwVer) *fwVer = pinfo->version;
		if(chipId) *chipId = pinfo->chip_module;
		if(service) *service = pinfo->hal_service;
		if(in) *in = pinfo->input_support;
		if(out) *out = pinfo->output_support;
		if(peripheral) *peripheral = pinfo->peripheral_en;
		return vatek_success;
	}
	return vatek_memfail;
}

