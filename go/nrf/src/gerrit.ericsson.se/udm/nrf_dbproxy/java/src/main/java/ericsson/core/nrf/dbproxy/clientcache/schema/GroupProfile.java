package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;
import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;

public class GroupProfile
{
    private String group_profile_id;
    private HashMap  nf_type;
	private HashMap  group_id;
	private int      profile_type;
	private Long     supi_version;
    private ByteString raw_group_profile;
    public GroupProfile()
    {
        this.group_profile_id = "";
        this.nf_type = new HashMap();
        this.group_id = new HashMap();
        this.profile_type = Code.PROFILE_TYPE_EMPTY;
        this.supi_version = 0L;
        this.raw_group_profile = ByteString.EMPTY;
    }

    public String getGroupProfileID()
    {
        return group_profile_id;
    }

    public void setGroupProfileID(String group_profile_id)
    {
        this.group_profile_id = group_profile_id;
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
	
    public void setSupiVersion(Long supi_version)
    {
		this.supi_version = supi_version;
    }

    public Long getSupiVersion()
    {
		return this.supi_version;
    }
	
    public ByteString getData()
    {
        return raw_group_profile;
    }

    public void setData(ByteString raw_group_profile)
    {
        this.raw_group_profile = raw_group_profile;
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
        return "GroupProfile{" +
               "group_profile_id=" + group_profile_id +
               ", nf_type=" + nfTypeStr +
               ", group_id=" + groupIdStr +
               ", profile_type=" + profile_type +
               ", supi_version=" + supi_version.toString() +
               ", raw_group_profile=" + raw_group_profile.toString() +
               '}';
    }
}
