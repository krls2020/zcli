package vpn

func (h *Handler) dnsIsAlive() (bool, error) {
	return true, nil
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	//defer cancel()
	//
	//if !nettools.HasIPv6PingCommand() {
	//	return false, errors.New(i18n.VpnStatusDnsNoCheckFunction)
	//}
	//err := nettools.Ping(ctx, "node1.master.core.zerops")
	//if err != nil {
	//	h.logger.Error(err)
	//	return false, nil
	//}
	//return true, nil
}
