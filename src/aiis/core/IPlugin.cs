using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace core
{
    public interface IPlugin
    {
        string Name();

        string Desc();

        void Run();

        void Release();
    }

    public class PluginConfig
    {
        public static string DEFAULT_NAME = "plugin.Plugin";
        public static string FUNC_RELEASE = "Release";
        public static string FUNC_RUN = "Run";
        public static string FUNC_NAME = "Name";
        public static string FUNC_DESC = "Desc";
    }
}
