using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Nancy;

namespace aiis
{
    class PluginManager : NancyModule
    {
        public PluginManager()
        {
            // 
            Get["/"] = parameters => {
                return "ok";
            };
        }

        public void LoadFromConfig()
        {

        }
    }
}
