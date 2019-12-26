package zap_custom_logger

import "context"

func (e *ElasticPlugin) InsertEndpointLogObjectString(log []byte) error {
	_, err := e.Client.Index().
		Index(StringPtr(GetEndpointsLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).
		Type("_doc").
		BodyString(string(log)).
		Do(context.Background())
	if err != nil {
		return err
	}
	_, err = e.Client.Flush().Index(StringPtr(GetEndpointsLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticPlugin) InsertEndpointLogObjectStruct(log interface{}) error {
	_, err := e.Client.Index().
		Index(StringPtr(GetEndpointsLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).
		Type("_doc").
		BodyJson(log).
		Do(context.Background())
	if err != nil {
		return err
	}
	_, err = e.Client.Flush().Index(StringPtr(GetEndpointsLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticPlugin) InsertServiceLogObjectString(log []byte) error {
	_, err := e.Client.Index().
		Index(StringPtr(GetServiceLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).
		Type("_doc").
		BodyString(string(log)).
		Do(context.Background())
	if err != nil {
		return err
	}
	_, err = e.Client.Flush().Index(StringPtr(GetServiceLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticPlugin) InsertServiceLogObjectStruct(log interface{}) error {
	_, err := e.Client.Index().
		Index(StringPtr(GetServiceLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).
		Type("_doc").
		BodyJson(log).
		Do(context.Background())
	if err != nil {
		return err
	}
	_, err = e.Client.Flush().Index(StringPtr(GetServiceLogsIndexName(String(e.logger.GetConfigs().App.ServiceName)))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
