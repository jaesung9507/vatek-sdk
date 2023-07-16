#include <vatek_sdk_usbstream.h>

extern uint8_t* GetTsFrame(char* param);

typedef struct vatek_context {
	hvatek_devices hdevlist;
	hvatek_chip hchip;
	hvatek_usbstream hustream;
	hmux_core hmux;
	hmux_channel m_hchannel;
	usbstream_param usbcmd;
	uint8_t buf[CHIP_STREAM_SLICE_LEN];
} VatekContext;

vatek_result source_sync_get_buffer(void* param, uint8_t** pslicebuf) {
	VatekContext* ctx = (VatekContext*)param;
	uint8_t* buf = GetTsFrame((char*)ctx);
	memcpy(ctx->buf, buf, CHIP_STREAM_SLICE_LEN);
	free(buf);
	*pslicebuf = ctx->buf;
	return (vatek_result)1;
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
		if(ctx->hdevlist) {
			vatek_device_list_free(ctx->hdevlist);
		}

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

	return (char*)ctx;
}

int VatekUsbDeviceOpen(char* p) {
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
		if(pinfo) {
			if(status) *status = pinfo->status;
			if(fwVer) *fwVer = pinfo->version;
			if(chipId) *chipId = pinfo->chip_module;
			if(service) *service = pinfo->hal_service;
			if(in) *in = pinfo->input_support;
			if(out) *out = pinfo->output_support;
			if(peripheral) *peripheral = pinfo->peripheral_en;
			return vatek_success;
		}
	}
	return vatek_memfail;
}

int VatekUsbStreamOpen(char *p) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx) {
		return vatek_usbstream_open(ctx->hchip, &ctx->hustream);
	}
	return vatek_memfail;
}

int VatekUsbStreamStart(char* p) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx) {
		ctx->usbcmd.mode = ustream_mode_sync;
		ctx->usbcmd.sync.param = ctx;
		ctx->usbcmd.sync.getbuffer = source_sync_get_buffer;

		return vatek_usbstream_start(ctx->hustream, &ctx->usbcmd);
	}
	return vatek_memfail;
}

int GetVatekUsbStreamStatus(char* p, int* status, uint32_t* cur, uint32_t* data, uint32_t* mode) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx && ctx->hustream) {
		Ptransform_info pinfo = NULL;
		usbstream_status s = vatek_usbstream_get_status(ctx->hustream, &pinfo);
		if(status) *status = s;
		if(pinfo) {
			if(cur)	*cur = pinfo->info.cur_bitrate;
			if(data) *data = pinfo->info.data_bitrate;
			if(mode) *mode = pinfo->mode;
		}
		return vatek_success;
	}
	return vatek_memfail;
}
