import { ref } from 'vue'
import { defineStore } from "pinia";
export const useNavStore = defineStore('Nav', () => {
    const nav = ref('Reverse')
    const command = ref('')
    const payloadContent = ref('')
    const dialogVisiable = ref(false)
    function setNav(value) {
        nav.value = value
    }

    function setPayloadContent(value) {
        dialogVisiable.value = true
        payloadContent.value = value
    }
    return { nav, setNav, dialogVisiable, payloadContent, setPayloadContent }
});