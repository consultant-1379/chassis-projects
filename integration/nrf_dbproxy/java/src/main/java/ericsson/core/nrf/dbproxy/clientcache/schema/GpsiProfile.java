package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import java.util.HashMap;

public class GpsiProfile {

  private String gpsi_profile_id;
  private HashMap nf_type;
  private HashMap group_id;
  private int profile_type;
  private Long gpsi_version;
  private ByteString raw_gpsi_profile;

  public GpsiProfile() {
    this.gpsi_profile_id = "";
    this.nf_type = new HashMap();
    this.group_id = new HashMap();
    this.profile_type = Code.PROFILE_TYPE_EMPTY;
    this.gpsi_version = 0L;
    this.raw_gpsi_profile = ByteString.EMPTY;
  }

  public String getGpsiProfileID() {
    return gpsi_profile_id;
  }

  public void setGpsiProfileID(String gpsiProfileId) {
    this.gpsi_profile_id = gpsiProfileId;
  }


  public void addNFType(String nfType) {
    this.nf_type.put(nfType, 1);
  }

  public void addGroupID(String groupId) {
    this.group_id.put(groupId, 1);
  }

  public void setProfileType(int profileType) {
    this.profile_type = profileType;
  }

  public Long getGpsiVersion() {
    return this.gpsi_version;
  }

  public void setGpsiVersion(Long gpsiVersion) {
    this.gpsi_version = gpsiVersion;
  }

  public ByteString getData() {
    return raw_gpsi_profile;
  }

  public void setData(ByteString rawGpsiProfile) {
    this.raw_gpsi_profile = rawGpsiProfile;
  }

  public String toString() {
    String nfTypeStr = "";
    if (!nf_type.isEmpty()) {
      nfTypeStr = nf_type.toString();
    }
    String groupIdStr = "";
    if (!group_id.isEmpty()) {
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
