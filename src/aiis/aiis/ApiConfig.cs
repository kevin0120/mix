using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace aiis
{
    class ApiConfig
    {
        public static string API_PREFIX()
        {
            return "/" + API_VERSION;
        }
        
        public static string API_TAG()
        {
            return "api/";
        }

        private static string API_VERSION = "v1";
    }
}
