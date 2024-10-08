package ericsson.core.nrf.dbproxy.common;

public final class Code {

  public static final int REGISTERED_PROVISIONED = 0;
  public static final int REGISTERED_ONLY = 1;
  public static final int PROVISIONED_ONLY = 2;
  public static final int OPERATOR_LT = -2;
  public static final int OPERATOR_LE = -1;
  public static final int OPERATOR_EQ = 0;
  public static final int OPERATOR_GE = 1;
  public static final int OPERATOR_GT = 2;
  public static final int OPERATOR_REGEX = 3;
  public static final String NFPROFILE_INDICE = "ericsson-nrf-nfprofiles";
  public static final String NFHELPER_INDICE = "ericsson-nrf-nfhelper";
  public static final String SUBSCRIPTION_INDICE = "ericsson-nrf-subscriptions";
  public static final String NRFADDRESS_INDICE = "ericsson-nrf-nrfaddresses";
  public static final String GROUPPROFILE_INDICE = "ericsson-nrf-groupprofiles";
  public static final String IMSIPREFIXPROFILE_INDICE = "ericsson-nrf-imsiprefixprofiles";
  public static final String NRFPROFILE_INDICE = "ericsson-nrf-nrfprofiles";
  public static final String GPSIPROFILE_INDICE = "ericsson-nrf-gpsiprofiles";
  public static final String GPSIPREFIXPROFILE_INDICE = "ericsson-nrf-gpsiprefixprofiles";
  public static final String CACHENFPROFILE_INDICE = "ericsson-nrf-cachenfprofiles";
  public static final String DISTRIBUTEDLOCK_INDICE = "ericsson-nrf-distributedlock";
  public static final String REGIONNFINFO_INDICE = "ericsson-nrf-regionnfinfo";
  public static final int NOT_CHANGED = 0;
  public static final int KVDB_LOCATOR_NAME_CHANGED = 1;
  public static final int KVDB_LOCATOR_PORT_CHANGED = 2;
  public static final int KVDB_LOCATOR_IP_CHANGED = 3;
  public static final int KVDB_REGION_NAME_CHANGED = 4;
  public static final int KEY_MAX_LENGTH = 1024;
  public static final int IMSI_MAX_LENGTH = 15;
  public static final int FRAGMENT_BLOCK_NF_PROFILE = 20;
  public static final int FRAGMENT_BLOCK_GROUP_PROFILE = 50;
  public static final int FRAGMENT_BLOCK_GPSI_PROFILE = 51;
  public static final int FRAGMENT_SESSION_ACTIVE_TIME = 5000;
  public static final int SUCCESS = 2000;
  public static final int SUBSCRIPTION_MONITOR_SUCCESS = 2003;
  public static final int CREATED = 2001;
  public static final int DATA_NOT_EXIST = 2002;
  public static final int VALID = 0;
  public static final int EMPTY_NF_INSTANCE_ID = 4000;
  public static final int EMPTY_NRF_ADDRESS_ID = 4001;
  public static final int EMPTY_SUBSCRIPTION_ID = 4002;
  public static final int EMPTY_NF_PROFILE = 4006;
  public static final int NF_PROFILE_FORMAT_ERROR = 4007;
  public static final int NF_INSTANCE_ID_LENGTH_EXCEED_MAX = 4008;
  public static final int NRF_ADDRESS_ID_LENGTH_EXCEED_MAX = 4009;
  public static final int SUBSCRIPTION_ID_LENGTH_EXCEED_MAX = 4010;
  public static final int INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE = 4011;
  public static final int INVALID_NRF_ADDRESS_ID_NOT_SERIALIZABLE = 4012;
  public static final int INVALID_SUBSCRIPTION_ID_NOT_SERIALIZABLE = 4013;
  public static final int EMPTY_FRAGMENT_SESSION_ID = 4015;
  public static final int EMPTY_SUBSCRIPTION_FILTER = 4016;
  public static final int EMPTY_SUBSCRIPTION_DATA = 4017;
  public static final int EMPTY_NRF_ADDRESS_FILTER = 4018;
  public static final int EMPTY_NRF_ADDRESS_DATA = 4019;
  public static final int EMPTY_GROUP_PROFILE_ID = 4020;
  public static final int GROUP_PROFILE_ID_LENGTH_EXCEED_MAX = 4021;
  public static final int EMPTY_GROUP_PROFILE_FILTER = 4022;
  public static final int EMPTY_GROUP_PROFILE_DATA = 4023;
  public static final int INVALID_GROUP_PROFILE_ID_NOT_SERIALIZABLE = 4024;
  public static final int EMPTY_NRF_INSTANCE_ID = 4025;
  public static final int NRF_INSTANCE_ID_LENGTH_EXCEED_MAX = 4026;
  public static final int EMPTY_RAW_NRF_PROFILE = 4027;
  public static final int EMPTY_SEARCH_IMSI = 4028;
  public static final int EMPTY_SEARCH_GPSI = 4029;
  public static final int EMPTY_GPSI_PROFILE_ID = 4030;
  public static final int GPSI_PROFILE_ID_LENGTH_EXCEED_MAX = 4031;
  public static final int EMPTY_GPSI_PROFILE_FILTER = 4032;
  public static final int EMPTY_GPSI_PROFILE_DATA = 4033;
  public static final int INVALID_GPSI_PROFILE_ID_NOT_SERIALIZABLE = 4034;
  public static final int SEARCH_IMSI_LENGTH_EXCEED_MAX = 4035;
  public static final int SEARCH_GPSI_LENGTH_EXCEED_MAX = 4036;
  public static final int EMPTY_CACHE_NF_INSTANCE_ID = 4037;
  public static final int CACHE_NF_INSTANCE_ID_LENGTH_EXCEED_MAX = 4038;
  public static final int EMPTY_RAW_CACHE_NF_PROFILE = 4039;
  public static final int GROUP_PROFILE_INVALID_PROFILE_TYPE = 4040;
  public static final int GPSI_PROFILE_INVALID_PROFILE_TYPE = 4041;
  public static final int INVALID_RANGE = 4050;
  public static final int INVALID_PROVISIONED = 4051;
  public static final int ATTRIBUTE_NOT_KNOWN = 4052;
  public static final int EMPTY_ATTRIBUTE_NAME = 4056;
  public static final int SEARCH_ATTRIBUTE_MISSED = 4057;
  public static final int SEARCH_VALUE_MISSED = 4058;
  public static final int EMPTY_SEARCH_VALUE = 4059;
  public static final int INVALID_ATTRIBUTE_OPERATOR = 4060;
  public static final int EMPTY_SEARCH_EXPRESSION = 4061;
  public static final int EMPTY_AND_EXPRESSION = 4062;
  public static final int EMPTY_OR_EXPRESSION = 4063;
  public static final int EMPTY_META_EXPRESSION = 4064;
  public static final int EMPTY_NFPROFILE_FILTER = 4066;
  public static final int INVALID_PROV_VERSION = 4067;
  public static final int CACHE_NF_PROFILE_FORMAT_ERROR = 4068;
  public static final int NFMESSAGE_PROTOCOL_ERROR = 4100;
  public static final int INTERNAL_ERROR = 5000;
  public static final int PROFILE_TYPE_EMPTY = 0;
  public static final int PROFILE_TYPE_GROUPID = 1;
  public static final int PROFILE_TYPE_INSTANCEID = 2;
  public static final int BAD_REQUEST = 4000;
  public static final String NFTYPE_NRF = "NRF";
  public static final String NFTYPE_NRFINFO = "NRFINFO";
  public static final String NFTYPE_AMF = "AMF";
  public static final String NFTYPE_AUSF = "AUSF";
  public static final String NFTYPE_BSF = "BSF";
  public static final String NFTYPE_CHF = "CHF";
  public static final String NFTYPE_PCF = "PCF";
  public static final String NFTYPE_SMF = "SMF";
  public static final String NFTYPE_UDM = "UDM";
  public static final String NFTYPE_UDR = "UDR";
  public static final String NFTYPE_UPF = "UPF";
  public static final String PATCHINSTID = "INSTID";
  public static final int PATCH_ADD = 1;
  public static final int PATCH_REMOVE = 2;
  public static final int PATCH_REPLACE = 3;
  private Code() {
  }
}
