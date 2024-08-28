
import os.path
import json
import logging
import jsonpatch
import sys
from hyper import tls
sys.path.append('../../RobotLib/Libraries/')
import tls_nocheck
from http2_client import http2_client
from http2_message_handler import http2_message_handler


class generate_NRF_NST_data(object):

    def __init__(self):
        self.hc = http2_client()
        self.hhander = http2_message_handler()

    def json_replace(self, template_profile, target_key, target):
        self.hhander.json_load(template_profile)
        self.hhander.json_patch_generate("replace",target_key,target)
        template_json_file = self.hhander.json_generate()
        return template_json_file
   

    def jsonfile_prepare(self, nf_type, host=None, port=None):
        if( nf_type == 'UDR' or nf_type == 'all'):
            groudId = 101
            ft_groupid_snssais_udr = open("nrf_disc_h2load_udr_groupid_snssais_queryURI", 'w')
            for id in range (101,171):
                template_json_file = self.json_replace("./NST_NFRegister_udr.json","/nfInstanceId","12345678-9udr-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/udrInfo/groupId","gid"+str(groudId))
                if( id < 121 ):
                    template_json_file = self.json_replace(template_json_file,"/nfServices/0/serviceName","nudr-comm")
                    ft_groupid_snssais_udr.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=nudr-comm&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UDR&requester-nf-type=AMF&supi=imsi-4600"+str(groudId)+"&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%220%22%7D&snssais=%7B%22sst%22%3A+1,%22sd%22%3A+%221%22%7D"+'\n')                   
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9udr-def0-1000-100000000"+str(id),template_json_file)
                self.hc.http2_client_receive()
                if id % 2 == 1:                   
                    template_provision_json_file = self.json_replace("./imsi_template_Profile.json","/supiRanges/0/pattern","^imsi-4600"+str(groudId))                  
                    template_provision_json_file = self.json_replace(template_provision_json_file,"/groupId","gid"+str(groudId))
                    print(template_provision_json_file)
                    self.hc.http2_client_send("POST","/nnrf-prov/v1/group-profile",template_provision_json_file)
                    self.hc.http2_client_receive()
                    groudId = groudId + 1                                   
                print ("groudId:%d, id:%d", groudId, id)                
            ft_groupid_snssais_udr.close()

        if(nf_type == 'UDM' or nf_type == 'all'):
            supi = 101                   
            ft_supi_snssais_udm = open("nrf_disc_h2load_udm_supi_snssais_queryURI", 'w')
            for id in range (101,191):
                template_json_file = self.json_replace("./NST_NFRegister_udm.json","/nfInstanceId","12345678-9udm-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/udmInfo/supiRanges/0/pattern","^imsi-600"+str(supi)+"\\d{4}$")               
                if( id < 121 ):
                    template_json_file = self.json_replace(template_json_file,"/nfServices/0/serviceName","nudm-test")
                    template_json_file = self.json_replace(template_json_file,"/nfServices/1/serviceName","nudm-test100")       
                    ft_supi_snssais_udm.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=nudm-test&service-names=nudm-test100&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UDM&requester-nf-type=UDR&supi=imsi-600"+str(supi)+"9999"+"&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%220%22%7D&snssais=%7B%22sst%22%3A+1,%22sd%22%3A+%221%22%7D"+'\n')
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9udm-def0-1000-100000000"+str(id),template_json_file)
                self.hc.http2_client_receive()
                if id % 3 == 2:
                    supi = supi + 1
            ft_supi_snssais_udm.close()
               

        if(nf_type == 'AUSF' or nf_type == 'all'):
            supi = 101
            for id in range (101,191):
                template_json_file = self.json_replace("./NST_NFRegister_ausf.json","/nfInstanceId","12345678-ausf-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/ausfInfo/supiRanges/0/pattern","^nai-"+str(supi)+".+@company\\.com$")
                if( id < 121 ):
                    template_json_file = self.json_replace(template_json_file,"/nfServices/0/serviceName","nausf-test")                  
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-ausf-def0-1000-100000000"+str(id),template_json_file)
                self.hc.http2_client_receive()
                if id % 3 == 2:
                    supi = supi + 1                    
                print ("supi:%d, id:%d", supi, id)


        if(nf_type == 'AMF' or nf_type == 'all'):
            tai = 101
            amfSetId = 101
            snssais = 101
            nsilist = 101            
            ft_tai_nsilist_amf = open("nrf_disc_h2load_amf_tai_nsilist_queryURI", 'w')
            for cycle in range(6,8):
                for id in range (101,501):
                    template_json_file = self.json_replace("./NST_NFRegister_amf.json","/nfInstanceId","12345678-9amf-def0-1000-10000000"+str(cycle)+str(id))
                    template_json_file = self.json_replace(template_json_file,"/amfInfo/amfSetId","5012"+str(amfSetId))
                    template_json_file = self.json_replace(template_json_file,"/amfInfo/taiList/0/tac","Abc"+str(tai))
                    template_json_file = self.json_replace(template_json_file,"/nsiList/0","111"+str(nsilist))
                    template_json_file = self.json_replace(template_json_file,"/sNssais/0/sd","AAA"+str(snssais))
                    ft_tai_nsilist_amf.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=namf-comm&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=AMF&requester-nf-type=UDM&tai=%7B%22plmnId%22%3A+%7B%22mcc%22%3A+%22310%22%2C%22mnc%22%3A+%22010%22%7D%2C%22tac%22%3A+%22Abc"+str(tai)+"%22%7D&nsi-list=111"+str(nsilist)+'\n')                   
                    print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                    print(template_json_file)
                    self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9amf-def0-1000-10000000"+str(cycle)+str(id),template_json_file)
                    self.hc.http2_client_receive()
                    if id % 2 == 1:
                        tai = tai + 1
                    if id % 16 == 15:       
                        amfSetId = amfSetId + 1
                    if id % 30 == 29:       
                        nsilist = nsilist + 1
                    if id % 90 == 89:       
                        snssais = snssais + 1
                    print ("tai:%d, amfSetId:%d, nsilist:%d, snssais:%d, id:%d, cycle:%d", tai, amfSetId,  nsilist, snssais, id, cycle)    
            ft_tai_nsilist_amf.close()


        if(nf_type == 'SMF' or nf_type == 'all'):
            tai = 101
            dnn = 101
            snssais = 101
            nsilist = 101           
            ft_tai_nsilist_smf = open("nrf_disc_h2load_smf_tai_nsilist_queryURI", 'w')
            ft_instancd_id_smf = open("nrf_disc_h2load_smf_instance_id_queryURI", 'w')
            for id in range (101,701):
                template_json_file = self.json_replace("./NST_NFRegister_smf.json","/nfInstanceId","12345678-9smf-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/smfInfo/taiList/0/tac","Abc"+str(tai))
                template_json_file = self.json_replace(template_json_file,"/smfInfo/dnnList/0","province"+str(dnn))
                template_json_file = self.json_replace(template_json_file,"/nsiList/0","111"+str(nsilist))
                template_json_file = self.json_replace(template_json_file,"/sNssais/0/sd","AAA"+str(snssais))                
                ft_tai_nsilist_smf.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=nsmf-test&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=SMF&requester-nf-type=UDM&tai=%7B%22plmnId%22%3A+%7B%22mcc%22%3A+%22310%22%2C%22mnc%22%3A+%22010%22%7D%2C%22tac%22%3A+%22Abc"+str(tai)+"%22%7D&nsi-list=111"+str(nsilist)+'\n')               
                ft_instancd_id_smf.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=nsmf-test&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=SMF&requester-nf-type=UDM&target-nf-instance-id=12345678-9smf-def0-1000-100000000"+str(id)+'\n')
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9smf-def0-1000-100000000"+str(id),template_json_file)
                self.hc.http2_client_receive()
                if id % 2 == 1:
                    tai = tai + 1                   
                if id % 20 == 19:       
                    dnn = dnn + 1
                if id % 22 == 21:       
                    nsilist = nsilist + 1
                if id % 66 == 65:       
                    snssais = snssais + 1
                print ("tai:%d, dnn:%d, nsilist:%d, snssais:%d, id:%d", tai, dnn, nsilist, snssais, id)
            ft_tai_nsilist_smf.close()
            ft_instancd_id_smf.close()

        if(nf_type == 'UPF' or nf_type == 'all'):
            dnn = 101
            nsilist = 101
            snssais = 101
            ft_dnn_snssais_upf = open("nrf_disc_h2load_upf_dnn_snssais_queryURI", 'w')           
            for id in range (101,701):
                template_json_file = self.json_replace("./NST_NFRegister_upf.json","/nfInstanceId","12345678-9upf-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/upfInfo/sNssaiUpfInfoList/0/dnnUpfInfoList/0/dnn","province"+str(dnn))
                template_json_file = self.json_replace(template_json_file,"/nsiList/0","111"+str(nsilist))
                template_json_file = self.json_replace(template_json_file,"/sNssais/0/sd","AAA"+str(snssais))
                template_json_file = self.json_replace(template_json_file,"/upfInfo/sNssaiUpfInfoList/0/sNssai/sd","AAA"+str(snssais))
                ft_dnn_snssais_upf.write("http://"+str(host)+":"+str(port)+"/nnrf-disc/v1/nf-instances?service-names=nupf-test&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UPF&requester-nf-type=UDM&dnn=province"+str(dnn)+"&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%22AAA"+str(snssais)+"%22%7D"+'\n')               
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9upf-def0-1000-100000000"+str(id),template_json_file)     
                self.hc.http2_client_receive()
                if id % 20 == 19:       
                    dnn = dnn + 1
                if id % 22 == 21:       
                    nsilist = nsilist + 1
                if id % 60 == 59:       
                    snssais = snssais + 1                   
                print ("dnn:%d, nsilist:%d, snssais:%d, id:%d", dnn, nsilist, snssais, id)
            ft_dnn_snssais_upf.close()


        if(nf_type == 'PCF' or nf_type == 'all'):
            dnn = 101           
            for id in range (101,191):
                template_json_file = self.json_replace("./NST_NFRegister_pcf.json","/nfInstanceId","12345678-9pcf-def0-1000-100000000"+str(id))
                template_json_file = self.json_replace(template_json_file,"/pcfInfo/dnnList/0","province"+str(dnn)+".mnc012.mcc345.gprs")
                if( id < 121 ):
                    template_json_file = self.json_replace(template_json_file,"/nfServices/0/serviceName","npcf-test")                                     
                print("!!!!!!!!!!!!!!!!!!!!!!template json file")
                print(template_json_file)
                self.hc.http2_client_send("PUT","/nnrf-nfm/v1/nf-instances/12345678-9pcf-def0-1000-100000000"+str(id),template_json_file)     
                self.hc.http2_client_receive()
                if id % 3 == 2:       
                    dnn = dnn + 1                   
                print ("dnn:%d, id:%d", dnn, id)



if __name__ == '__main__':
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    if len(sys.argv) ==2 and sys.argv[1] == "-h":
        print("[Usage]: python generate_NRF_NST_data.py [eric-nrf-management ClusterIP] [port] [NF Type]")
        print("[NF Type]: AMF, UDR, UDM, AUSF, PCF, SMF, UPF, all")
    else:
        if len(sys.argv) ==4:
            host = sys.argv[1]
            port = sys.argv[2]
            nf_type = sys.argv[3]
    

        generate_data = generate_NRF_NST_data()
        generate_data.hc.http2_client_connect(host, port)
        generate_data.jsonfile_prepare(nf_type, host, port)
        generate_data.hc.http2_client_disconnect()
