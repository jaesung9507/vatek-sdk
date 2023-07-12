#include <stdio.h>
#include <vatek_sdk_usbstream.h>
#include <core/ui/ui_props/ui_props_chip.h>

#define _disp_l(fmt,...)	printf("	"fmt"\r\n",##__VA_ARGS__)
#define _disp_err(fmt,...)	printf("	error - "fmt"\r\n",##__VA_ARGS__)

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
} VatekContext;

static usbstream_param usbcmd = {
	.mode = ustream_mode_sync,
	.remux = ustream_remux_pcr,
	.pcradjust = pcr_adjust,
	.r2param.freqkhz = 473000,							/* output _rf frequency */
	.r2param.mode = r2_cntl_path_0,
	.modulator = {
		6,
		modulator_dvb_t,
		ifmode_disable,0,0,
		.mod = {.dvb_t = {dvb_t_qam64, fft_8k, guard_interval_1_16, coderate_5_6},},
	},
	.sync = {NULL,NULL},
};

void printf_chip_info(Pchip_info pchip) {
	_disp_l("-------------------------------------");
	_disp_l("	   Chip Information");
	_disp_l("-------------------------------------");
	_disp_l("%-20s : %s", "Chip Status", ui_enum_get_str(chip_status, pchip->status));
	_disp_l("%-20s : %08x", "FW Version", pchip->version);
	_disp_l("%-20s : %08x", "Chip  ID", pchip->chip_module);
	_disp_l("%-20s : %08x", "Service", pchip->hal_service);
	_disp_l("%-20s : %08x", "Input", pchip->input_support);
	_disp_l("%-20s : %08x", "Output", pchip->output_support);
	_disp_l("%-20s : %08x", "Peripheral", pchip->peripheral_en);
	_disp_l("%-20s : %d", "SDK Version", vatek_version());
	_disp_l("=====================================\r\n");
}

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

void FreeVatekContext(char* p) {
	VatekContext* ctx = (VatekContext*)p;
	if(ctx) {
		if(ctx->hustream) {
			vatek_result nres = vatek_usbstream_stop(ctx->hustream);
			if(!is_vatek_success(nres))
				_disp_err("stop usb_stream fail : %d", nres);
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
}

char* NewVatekContext() {
	VatekContext* ctx = (VatekContext*)malloc(sizeof(VatekContext));
	memset(ctx, 0, sizeof(VatekContext));
	modulator_param_reset(modulator_atsc, &usbcmd.modulator);
	ctx->streamsource.hsource = NULL;
	ctx->streamsource.start = ts_stream_start;
	ctx->streamsource.stop = ts_stream_stop;
	ctx->streamsource.get = ts_stream_get;
	ctx->streamsource.check = ts_stream_check;
	ctx->streamsource.free = ts_stream_free;

	vatek_result nres = vatek_device_list_enum(DEVICE_BUS_USB, service_transform, &ctx->hdevlist);
	if(is_vatek_success(nres)) {
		if(nres == 0) {
			nres = vatek_nodevice;
			_disp_err("can not found device.");
		} else {
			nres = vatek_device_open(ctx->hdevlist, 0, &ctx->hchip);
			if(!is_vatek_success(nres)) {
				_disp_err("open device fail : %d", nres);
			} else {
				Pchip_info pinfo = vatek_device_get_info(ctx->hchip);
				printf_chip_info(pinfo);
				nres = ctx->streamsource.start(ctx->streamsource.hsource);
			}
		}
	} else {
		_disp_err("enum device fail : %d",nres);
	}

	if(!is_vatek_success(nres)) {
		FreeVatekContext((char*)ctx);
		return NULL;
	}
	return (char*)ctx;
}
