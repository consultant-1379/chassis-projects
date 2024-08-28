package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;
import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;

public class GpsiProfile
{
    private String gpsi_profile_id;
    private HashMap  nf_type;
	private HashMap  group_id;
	private int      profile_type;
	private Long      gpsi_version;
    private ByteString raw_gpsi_profile;
    public GpsiProfile()
    {
        this.gpsi_profile_id = "";
        this.nf_type = new HashMap();
		this.group_id = new HashMap();
        this.profile_type = Code.PROFILE_TYPE_EMPTY;
        this.gpsi_version = 0L;
        this.raw_gpsi_profile = ByteString.EMPTY;
    }

    public String getGpsiProfileID()
    {
        return gpsi_profile_id;
    }

    public void setGpsiProfileID(String gpsi_profile_id)
    {
        this.gpsi_profile_id = gpsi_profile_id;
    }


    public void addNFType(String nf_type)
    {
        this.nf_type.put(nf_type, 1);
    }
	
    public void addGroupID(String group_id)
    {
		this.group_id.put(group_id, 1);
    }
	
    public void setProfileType(int profile_type)
    {
		this.profile_type = profile_type;
    }
	
    public void setGpsiVersion(Long gpsi_version)
    {
		this.gpsi_version = gpsi_version;
    }

    public Long getGpsiVersion()
    {
		return this.gpsi_version;
    }

    public ByteString getData()
    {
        return raw_gpsi_profile;
    }

    public void setData(ByteString raw_gpsi_profile)
    {
        this.raw_gpsi_profile = raw_gpsi_profile;
    }

    public String toString()
    {
		String nfTypeStr = "";
		if (nf_type.isEmpty() == false) {
			nfTypeStr = nf_type.toString();
		}
		String groupIdStr = "";
		if (group_id.isEmpty() == false) {
			groupIdStr = group_id.toString();
		}
        return "GpsiProfile{" +
               "gpsi_profile_id=" + gpsi_profile_id +
               ", nf_type=" + nfTypeStr +
               ", group_id=" + groupIdStr +
               ", profile_type=" + profile_type +
               ", gpsi_version=" + gpsi_version +
               ", raw_gpsi_profile=" + raw_gpsi_profile.toString() +
               '}';
    }
}
